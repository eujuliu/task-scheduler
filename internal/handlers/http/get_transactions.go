package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"

	"github.com/gin-gonic/gin"
)

type GetTransactionsHandler struct {
	getTransactionsService *services.GetTransactionsService
}

func NewGetTransactionsHandler(
	getTransactionsService *services.GetTransactionsService,
) *GetTransactionsHandler {
	return &GetTransactionsHandler{
		getTransactionsService: getTransactionsService,
	}
}

func (h *GetTransactionsHandler) Handle(c *gin.Context) {
	userId, ok := helpers.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This access token is not valid",
			"success": false,
		})
	}

	transactions := h.getTransactionsService.Execute(userId)
	result := []map[string]any{}

	for _, transaction := range transactions {
		result = append(result, map[string]any{
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

	c.JSON(http.StatusOK, result)
}
