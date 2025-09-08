package http_handlers

import (
	"net/http"
	"scheduler/internal/entities"
	stripe_paymentgateway "scheduler/internal/payment_gateway/stripe"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/postgres"

	"github.com/gin-gonic/gin"
)

type BuyCreditsRequest struct {
	Credits  int    `json:"credits"  binding:"required,gte=10,lte=100"`
	Currency string `json:"currency" binding:"required,iso4217"`
}

func BuyCredits(c *gin.Context) {
	var json BuyCreditsRequest

	idempotencyKey := c.Request.Header.Get("Idempotency-Key")

	if idempotencyKey == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error":   "Missing idempotency key header (Idempotency-Key)",
				"code":    http.StatusBadRequest,
				"success": false,
			},
		)

		return
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error(), "code": http.StatusBadRequest, "success": false},
		)
		return
	}

	postgres.DB.BeginTransaction()

	userRepository := postgres_repos.NewPostgresUserRepository()
	transactionRepository := postgres_repos.NewPostgresTransactionRepository()

	customerPaymentGateway := stripe_paymentgateway.NewStripeCustomerPaymentGateway()
	paymentPaymentGateway := stripe_paymentgateway.NewStripePaymentPaymentGateway()

	createTransactionService := services.NewCreateTransactionService(
		userRepository,
		transactionRepository,
		customerPaymentGateway,
		paymentPaymentGateway,
	)

	userId, ok := helpers.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This access token is not valid",
			"success": false,
		})

		return
	}

	transaction, err := createTransactionService.Execute(
		userId,
		json.Credits,
		json.Currency,
		entities.TypeTransactionPurchase,
		"",
		idempotencyKey,
	)
	if err != nil {
		_ = postgres.DB.RollbackTransaction()

		_ = c.Error(err)

		return
	}

	_ = postgres.DB.CommitTransaction()

	c.JSON(http.StatusOK, gin.H{
		"id":          transaction.GetId(),
		"credits":     transaction.GetCredits(),
		"amount":      transaction.GetAmount(),
		"currency":    transaction.GetCurrency(),
		"status":      transaction.GetStatus(),
		"createdAt":   transaction.GetCreatedAt(),
		"updateAt":    transaction.GetUpdatedAt(),
		"referenceId": transaction.GetReferenceId(),
	})
}
