package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionResponse struct {
	ID        string `json:"id"`
	Credits   int    `json:"credits"`
	Amount    int    `json:"amount"`
	Currency  string `json:"currency"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

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

// @Summary		Get transactions
// @Description	Get list of transactions for the user
// @Tags			transactions
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			offset	query		int		false	"Offset"
// @Param			limit	query		int		false	"Limit"
// @Param			orderBy	query		string	false	"Order by"
// @Success		200		{array}		TransactionResponse
// @Failure		404		{object}	errors.Error
// @Router			/transactions [get]
func (h *GetTransactionsHandler) Handle(c *gin.Context) {
	userId, _ := helpers.GetUserID(c)

	var offset int
	var limit int
	orderBy, _ := c.GetQuery("orderBy")

	if value, ok := c.GetQuery("offset"); ok {
		v, _ := strconv.Atoi(value)

		offset = v
	}

	if value, ok := c.GetQuery("limit"); ok {
		v, _ := strconv.Atoi(value)

		limit = v
	}

	transactions := h.getTransactionsService.Execute(userId, &offset, &limit, &orderBy)
	result := []TransactionResponse{}

	for _, transaction := range transactions {
		result = append(result, TransactionResponse{
			ID:        transaction.GetId(),
			Credits:   transaction.GetCredits(),
			Amount:    transaction.GetAmount(),
			Currency:  transaction.GetCurrency(),
			Status:    transaction.GetStatus(),
			Type:      transaction.GetType(),
			CreatedAt: transaction.GetCreatedAt().Format(time.RFC3339),
			UpdatedAt: transaction.GetUpdatedAt().Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, result)
}
