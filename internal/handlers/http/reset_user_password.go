package http_handlers

import (
	"net/http"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/postgres"

	"github.com/gin-gonic/gin"
)

type ResetUserPasswordRequest struct {
	TokenID     string `json:"tokenId"  binding:"required"`
	NewPassword string `json:"password" binding:"required"`
}

func ResetUserPassword(c *gin.Context) {
	var json ResetUserPasswordRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error(), "code": http.StatusBadRequest, "success": false},
		)
		return
	}

	userRepository := postgres_repos.NewPostgresUserRepository()
	passwordRepository := postgres_repos.NewPostgresPasswordRepository()
	resetPasswordService := services.NewResetUserPasswordService(userRepository, passwordRepository)

	postgres.DB.BeginTransaction()

	err := resetPasswordService.Execute(json.TokenID, json.NewPassword)
	if err != nil {
		_ = postgres.DB.RollbackTransaction()

		_ = c.Error(err)

		return
	}

	_ = postgres.DB.CommitTransaction()

	c.JSON(http.StatusOK, gin.H{
		"message": "Password has been reset!",
	})
}
