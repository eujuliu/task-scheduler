package http_handlers

import (
	"net/http"
	"scheduler/internal/entities"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/postgres"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BuyCreditsRequest struct {
	Credits  int    `json:"credits"  binding:"required,gte=10,lte=100"`
	Currency string `json:"currency" binding:"required,iso4217"`
}

type BuyCreditsResponse struct {
	ID          string `json:"id"`
	Credits     int    `json:"credits"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	ReferenceID string `json:"referenceId"`
}

type BuyCreditsHandler struct {
	db                       *postgres.Database
	createTransactionService *services.CreateTransactionService
}

func NewBuyCreditsHandler(
	db *postgres.Database,
	createTransactionService *services.CreateTransactionService,
) *BuyCreditsHandler {
	return &BuyCreditsHandler{
		db:                       db,
		createTransactionService: createTransactionService,
	}
}

// @Summary		Buy credits
// @Description	Buy credits for the user
// @Tags			credits
// @Accept			json
// @Produce		json
// @Param			Idempotency-Key	header		string				true	"Idempotency Key"
// @Param			request			body		BuyCreditsRequest	true	"Buy credits request"
// @Success		200				{object}	BuyCreditsResponse
// @Failure		400				{object}	errors.Error
// @Failure		404				{object}	errors.Error
// @Router			/buy-credits [post]
func (h *BuyCreditsHandler) Handle(c *gin.Context) {
	idempotencyKey := c.Request.Header.Get("Idempotency-Key")

	if idempotencyKey == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"id":    uuid.NewString(),
				"error": "Missing idempotency key header (Idempotency-Key)",
				"code":  http.StatusBadRequest,
			},
		)

		return
	}

	var json BuyCreditsRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error(), "code": http.StatusBadRequest, "success": false},
		)
		return
	}

	h.db.BeginTransaction()

	userId, _ := helpers.GetUserID(c)

	transaction, err := h.createTransactionService.Execute(
		userId,
		json.Credits,
		json.Currency,
		entities.TypeTransactionPurchase,
		"",
		idempotencyKey,
	)
	if err != nil {
		_ = h.db.RollbackTransaction()

		_ = c.Error(err)

		return
	}

	_ = h.db.CommitTransaction()

	response := BuyCreditsResponse{
		ID:          transaction.GetId(),
		Credits:     transaction.GetCredits(),
		Amount:      transaction.GetAmount(),
		Currency:    transaction.GetCurrency(),
		Status:      transaction.GetStatus(),
		CreatedAt:   transaction.GetCreatedAt().Format(time.RFC3339),
		UpdatedAt:   transaction.GetUpdatedAt().Format(time.RFC3339),
		ReferenceID: transaction.GetReferenceId(),
	}

	c.JSON(http.StatusOK, response)
}
