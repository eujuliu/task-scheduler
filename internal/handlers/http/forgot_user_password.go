package http_handlers

import (
	"net/http"
	"scheduler/internal/services"

	"github.com/gin-gonic/gin"
)

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type ForgotPasswordHandler struct {
	forgotPasswordService *services.ForgotUserPasswordService
}

func NewForgotPasswordHandler(
	forgotPasswordService *services.ForgotUserPasswordService,
) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		forgotPasswordService: forgotPasswordService,
	}
}

// @Summary		Forgot password
// @Description	Send recovery email for password reset
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		ForgotPasswordRequest	true	"Forgot password request"
// @Success		200		{object}	ForgotPasswordResponse
// @Failure		400		{object}	errors.Error
// @Failure		404		{object}	errors.Error
// @Router			/auth/forgot-password [post]
func (h *ForgotPasswordHandler) Handle(c *gin.Context) {
	var json ForgotPasswordRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error(), "code": http.StatusBadRequest, "success": false},
		)
		return
	}

	_, err := h.forgotPasswordService.Execute(json.Email)
	if err != nil {
		_ = c.Error(err)

		return
	}

	// Logic for send via email, use the service for create task with an email template

	response := ForgotPasswordResponse{
		Message: "Your recovery key was sent to your email",
	}

	c.JSON(http.StatusOK, response)
}
