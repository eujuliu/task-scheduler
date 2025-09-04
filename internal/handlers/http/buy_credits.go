package http_handlers

import (
	"log/slog"
	"net/http"
	"scheduler/internal/entities"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/postgres"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BuyCreditsRequest struct {
	Credits int `json:"credits" binding:"required"`
}

func BuyCredits(c *gin.Context) {
	var json BuyCreditsRequest

	idempotencyKey := c.Request.Header.Get("Idempotency-Key")

	slog.Info(idempotencyKey)

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

	// Payment gateway logic

	userRepository := postgres_repos.NewPostgresUserRepository()
	transactionRepository := postgres_repos.NewPostgresTransactionRepository()
	createTransactionService := services.NewCreateTransactionService(
		userRepository,
		transactionRepository,
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
		10,
		"BRL",
		entities.TypeTransactionPurchase,
		uuid.NewString(),
		idempotencyKey,
	)
	if err != nil {
		_ = postgres.DB.RollbackTransaction()

		_ = c.Error(err)

		return
	}

	_ = postgres.DB.CommitTransaction()

	c.JSON(http.StatusOK, gin.H{
		"id":        transaction.GetId(),
		"credits":   transaction.GetCredits(),
		"amount":    transaction.GetAmount(),
		"currency":  transaction.GetCurrency(),
		"status":    transaction.GetStatus(),
		"createdAt": transaction.GetCreatedAt(),
		"updateAt":  transaction.GetUpdatedAt(),
	})
}
