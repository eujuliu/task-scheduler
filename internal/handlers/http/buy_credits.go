package http_handlers

import (
	"net/http"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/postgres"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BuyCreditsRequest struct {
	Credits       int    `json:"credits"       binding:"required"`
	IdempotentKey string `json:"idempotentKey" binding:"required"`
}

func BuyCredits(c *gin.Context) {
	var json BuyCreditsRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		json.IdempotentKey,
	)
	if err != nil {
		_ = postgres.DB.RollbackTransaction()

		if e := errors.GetError(err); e != nil {
			c.JSON(e.Code, gin.H{
				"code":    e.Code,
				"message": e.Msg(),
			})
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Internal Server Error", "message": "contact the admin"},
		)
		return
	}

	_ = postgres.DB.CommitTransaction()

	c.JSON(http.StatusOK, gin.H{
		"id":            transaction.GetId(),
		"credits":       transaction.GetCredits(),
		"amount":        transaction.GetAmount(),
		"currency":      transaction.GetCurrency(),
		"status":        transaction.GetStatus(),
		"idempotentKey": transaction.GetIdempotencyKey(),
		"createdAt":     transaction.GetCreatedAt(),
		"updateAt":      transaction.GetUpdatedAt(),
	})
}
