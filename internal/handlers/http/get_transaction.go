package http_handlers

import (
	"net/http"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"

	"github.com/gin-gonic/gin"
)

func GetTransaction(c *gin.Context) {
	id := c.Param("id")

	userRepository := postgres_repos.NewPostgresUserRepository()
	transactionRepository := postgres_repos.NewPostgresTransactionRepository()
	getTransactionService := services.NewGetTransactionService(
		userRepository,
		transactionRepository,
	)

	userId, ok := helpers.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This access token is not valid",
		})
	}

	transaction, err := getTransactionService.Execute(userId, id)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        transaction.GetId(),
		"credits":   transaction.GetCredits(),
		"amount":    transaction.GetAmount(),
		"currency":  transaction.GetCurrency(),
		"status":    transaction.GetStatus(),
		"type":      transaction.GetType(),
		"createdAt": transaction.GetCreatedAt(),
		"updatedAt": transaction.GetUpdatedAt(),
	})
}
