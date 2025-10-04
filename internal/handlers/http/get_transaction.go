package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"time"

	"github.com/gin-gonic/gin"
)

type GetTransactionResponse struct {
	ID        string `json:"id"`
	Credits   int    `json:"credits"`
	Amount    int    `json:"amount"`
	Currency  string `json:"currency"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

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

// @Summary		Get transaction
// @Description	Get a specific transaction by ID
// @Tags			transactions
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"Transaction ID"
// @Success		200	{object}	GetTransactionResponse
// @Failure		404	{object}	errors.Error
// @Router			/transaction/{id} [get]
func (h *GetTransactionHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	userId, _ := helpers.GetUserID(c)

	transaction, err := h.getTransactionService.Execute(userId, id)
	if err != nil {
		_ = c.Error(err)

		return
	}

	response := GetTransactionResponse{
		ID:        transaction.GetId(),
		Credits:   transaction.GetCredits(),
		Amount:    transaction.GetAmount(),
		Currency:  transaction.GetCurrency(),
		Status:    transaction.GetStatus(),
		Type:      transaction.GetType(),
		CreatedAt: transaction.GetCreatedAt().Format(time.RFC3339),
		UpdatedAt: transaction.GetUpdatedAt().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}
