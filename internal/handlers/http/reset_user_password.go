package http_handlers

import (
	"net/http"
	"scheduler/internal/errors"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepository := postgres_repos.NewPostgresUserRepository()
	passwordRepository := postgres_repos.NewPostgresPasswordRepository()
	resetPasswordService := services.NewResetUserPasswordService(userRepository, passwordRepository)

	postgres.DB.BeginTransaction()

	err := resetPasswordService.Execute(json.TokenID, json.NewPassword)
	if err != nil {
		_ = postgres.DB.RollbackTransaction()

		if e := errors.GetError(err); e != nil {
			c.JSON(e.Code, gin.H{
				"code":    e.Code,
				"message": e.Msg(),
			})
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Internal Server Error", "message": "contact the admin"},
		)
		return
	}

	_ = postgres.DB.CommitTransaction()

	c.JSON(http.StatusOK, gin.H{
		"message": "Password has been reset!",
	})
}
