package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"

	"github.com/gin-gonic/gin"
)

type GetTransactionHandler struct {
	getTransactionService *services.GetTransactionService
}

func NewGetTransactionHandler(
	getTransactionService *services.GetTransactionService,
) *GetTransactionHandler {
	return &GetTransactionHandler{
		getTransactionService: getTransactionService,
	}
}

func (h *GetTransactionHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	userId, ok := helpers.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This access token is not valid",
		})
	}

	transaction, err := h.getTransactionService.Execute(userId, id)
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
