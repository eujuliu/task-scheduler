package http_handlers

import (
	"net/http"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"

	"github.com/gin-gonic/gin"
)

func Transactions(c *gin.Context) {
	userRepository := postgres_repos.NewPostgresUserRepository()
	transactionRepository := postgres_repos.NewPostgresTransactionRepository()
	getTransactionsService := services.NewGetTransactionsService(
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

	transactions := getTransactionsService.Execute(userId)
	result := []map[string]any{}

	for _, transaction := range transactions {
		result = append(result, map[string]any{
			"id":            transaction.GetId(),
			"credits":       transaction.GetCredits(),
			"amount":        transaction.GetAmount(),
			"currency":      transaction.GetCurrency(),
			"status":        transaction.GetStatus(),
			"type":          transaction.GetType(),
			"idempotentKey": transaction.GetIdempotencyKey(),
			"createdAt":     transaction.GetCreatedAt(),
			"updatedAt":     transaction.GetUpdatedAt(),
		})
	}

	c.JSON(http.StatusOK, result)
}
