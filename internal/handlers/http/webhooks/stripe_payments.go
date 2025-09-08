package http_webhooks

import (
	"encoding/json"
	"io"
	"net/http"
	"scheduler/internal/config"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

func StripePayments(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.Error(err)
		c.Status(http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEvent(
		body,
		c.GetHeader("Stripe-Signature"),
		config.Data.Stripe.EndpointSecret,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userRepository := postgres_repos.NewPostgresUserRepository()
	transactionRepository := postgres_repos.NewPostgresTransactionRepository()
	errorRepository := postgres_repos.NewPostgresErrorRepository()

	updatePurchaseService := services.NewUpdatePurchaseTransactionService(
		userRepository,
		transactionRepository,
		errorRepository,
	)

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent

		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			_ = c.Error(err)
			return
		}

		_, err := updatePurchaseService.Complete(paymentIntent.Metadata["transactionID"])
		if err != nil {
			_ = c.Error(err)

			return
		}

	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent

		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			_ = c.Error(err)
			return
		}

		_, err := updatePurchaseService.Fail(
			paymentIntent.Metadata["transactionID"],
			paymentIntent.LastPaymentError.Msg,
		)
		if err != nil {
			_ = c.Error(err)

			return
		}
	}

	c.Status(http.StatusOK)
}
