package http_webhooks

import (
	"encoding/json"
	"io"
	"net/http"
	"scheduler/internal/config"
	"scheduler/internal/interfaces"
	"scheduler/internal/queue"
	"scheduler/internal/services"
	"scheduler/pkg/postgres"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type StripePaymentUpdateWebhook struct {
	db                       *postgres.Database
	config                   *config.StripeConfig
	queue                    interfaces.IQueue
	updateTransactionService *services.UpdatePurchaseTransactionService
}

func NewStripePaymentUpdateWebhook(
	db *postgres.Database,
	config *config.StripeConfig,
	queue interfaces.IQueue,
	updateTransactionService *services.UpdatePurchaseTransactionService,
) *StripePaymentUpdateWebhook {
	return &StripePaymentUpdateWebhook{
		db:                       db,
		config:                   config,
		queue:                    queue,
		updateTransactionService: updateTransactionService,
	}
}

func (wh *StripePaymentUpdateWebhook) Hook(c *gin.Context) {
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
		wh.config.EndpointSecret,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent

		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			_ = c.Error(err)
			return
		}

		wh.db.BeginTransaction()

		transaction, err := wh.updateTransactionService.Complete(
			paymentIntent.Metadata["transactionID"],
		)
		if err != nil {
			_ = wh.db.RollbackTransaction()
			_ = c.Error(err)
			return
		}

		_ = wh.db.CommitTransaction()

		event, err := queue.NewEvent(transaction.GetUserId(), "success", map[string]any{
			"transactionId": transaction.GetId(),
			"when":          transaction.GetUpdatedAt(),
			"credits":       transaction.GetCreatedAt(),
			"cost":          transaction.GetAmount(),
			"currency":      transaction.GetCurrency(),
		})
		if err != nil {
			_ = c.Error(err)
			return
		}

		data, err := json.Marshal(event)
		if err != nil {
			_ = c.Error(err)
			return
		}

		err = wh.queue.Publish(queue.SEND_EVENTS_KEY, queue.EVENTS_EXCHANGE, data, event.ClientID)
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

		wh.db.BeginTransaction()
		transaction, err := wh.updateTransactionService.Fail(
			paymentIntent.Metadata["transactionID"],
			paymentIntent.LastPaymentError.Msg,
		)
		if err != nil {
			_ = wh.db.RollbackTransaction()
			_ = c.Error(err)
			return
		}

		_ = wh.db.CommitTransaction()

		event, err := queue.NewEvent(transaction.GetUserId(), "error", map[string]any{
			"transactionId": transaction.GetId(),
			"message":       paymentIntent.LastPaymentError.Msg,
			"when":          transaction.GetUpdatedAt(),
			"credits":       transaction.GetCreatedAt(),
			"cost":          transaction.GetAmount(),
			"currency":      transaction.GetCurrency(),
		})
		if err != nil {
			_ = c.Error(err)
			return
		}

		data, err := json.Marshal(event)
		if err != nil {
			_ = c.Error(err)
			return
		}

		err = wh.queue.Publish(queue.SEND_EVENTS_KEY, queue.EVENTS_EXCHANGE, data, event.ClientID)
		if err != nil {
			_ = c.Error(err)
			return
		}
	}

	c.Status(http.StatusOK)
}
