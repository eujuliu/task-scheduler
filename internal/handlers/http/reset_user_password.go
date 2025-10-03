package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/postgres"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResetUserPasswordRequest struct {
	TokenID     string `json:"tokenId"  binding:"required"`
	NewPassword string `json:"password" binding:"required"`
}

type ResetUserPasswordResponse struct {
	Message string `json:"message"`
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

// @Summary		Reset user password
// @Description	Reset password using token ID and new password
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		ResetUserPasswordRequest	true	"Reset password request"
// @Success		200		{object}	ResetUserPasswordResponse
// @Failure		400		{object}	errors.Error
// @Failure		404		{object}	errors.Error
// @Router			/auth/reset-password [post]
func (h *ResetUserPasswordHandler) Handle(c *gin.Context) {
	var json ResetUserPasswordRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"id": uuid.NewString(), "error": err.Error(), "code": http.StatusBadRequest},
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

	response := ResetUserPasswordResponse{
		Message: "Password has been reset!",
	}

	c.JSON(http.StatusOK, response)
}
