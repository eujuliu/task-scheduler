package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/postgres"

	"github.com/gin-gonic/gin"
)

type ResetUserPasswordRequest struct {
	TokenID     string `json:"tokenId"  binding:"required"`
	NewPassword string `json:"password" binding:"required"`
}

type ResetUserPasswordHandler struct {
	db                       *postgres.Database
	resetUserPasswordService *services.ResetUserPasswordService
}

func NewResetUserPasswordHandler(
	db *postgres.Database,
	resetUserPasswordService *services.ResetUserPasswordService,
) *ResetUserPasswordHandler {
	return &ResetUserPasswordHandler{
		db:                       db,
		resetUserPasswordService: resetUserPasswordService,
	}
}

func (h *ResetUserPasswordHandler) Handle(c *gin.Context) {
	var json ResetUserPasswordRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error(), "code": http.StatusBadRequest, "success": false},
		)
		return
	}

	h.db.BeginTransaction()

	err := h.resetUserPasswordService.Execute(json.TokenID, json.NewPassword)
	if err != nil {
		_ = h.db.RollbackTransaction()

		_ = c.Error(err)

		return
	}

	_ = h.db.CommitTransaction()

	c.JSON(http.StatusOK, gin.H{
		"message": "Password has been reset!",
	})
}
