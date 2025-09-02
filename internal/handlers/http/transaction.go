package http_handlers

import (
	"net/http"
	"scheduler/internal/errors"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"

	"github.com/gin-gonic/gin"
)

func Transaction(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
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
