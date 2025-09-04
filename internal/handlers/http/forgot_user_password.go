package http_handlers

import (
	"net/http"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"

	"github.com/gin-gonic/gin"
)

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func ForgotPassword(c *gin.Context) {
	var json ForgotPasswordRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error(), "code": http.StatusBadRequest, "success": false},
		)
		return
	}

	userRepository := postgres_repos.NewPostgresUserRepository()
	passwordRepository := postgres_repos.NewPostgresPasswordRepository()
	forgotPasswordService := services.NewForgotUserPasswordService(
		userRepository,
		passwordRepository,
	)

	_, err := forgotPasswordService.Execute(json.Email)
	if err != nil {
		_ = c.Error(err)

		return
	}

	// Logic for send via email, use the service for create task with an email template

	c.JSON(http.StatusOK, gin.H{
		"message": "Your recovery key was sent to your email",
	})
}
