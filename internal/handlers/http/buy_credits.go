package http_handlers

import (
	"net/http"
	"scheduler/internal/entities"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/postgres"

	"github.com/gin-gonic/gin"
)

type BuyCreditsRequest struct {
	Credits  int    `json:"credits"  binding:"required,gte=10,lte=100"`
	Currency string `json:"currency" binding:"required,iso4217"`
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

func (h *BuyCreditsHandler) Handle(c *gin.Context) {
	var json BuyCreditsRequest

	idempotencyKey := c.Request.Header.Get("Idempotency-Key")

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

	h.db.BeginTransaction()

	userId, ok := helpers.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This access token is not valid",
			"success": false,
		})

		return
	}

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

	c.JSON(http.StatusOK, gin.H{
		"id":          transaction.GetId(),
		"credits":     transaction.GetCredits(),
		"amount":      transaction.GetAmount(),
		"currency":    transaction.GetCurrency(),
		"status":      transaction.GetStatus(),
		"createdAt":   transaction.GetCreatedAt(),
		"updateAt":    transaction.GetUpdatedAt(),
		"referenceId": transaction.GetReferenceId(),
	})
}
