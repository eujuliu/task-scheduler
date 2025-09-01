package http_handlers

import (
	"net/http"
	"scheduler/internal/config"
	"scheduler/internal/errors"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var json RegisterRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conf := config.Instance
	userRepository := postgres_repos.NewPostgresUserRepository()
	createUserService := services.NewCreateUserService(userRepository)

	user, err := createUserService.Execute(json.Username, json.Email, json.Password)
	if err != nil {
		err, ok := err.(*errors.Error)

		if ok {
			c.JSON(err.Code, gin.H{
				"code":    err.Code,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	accessToken, err := utils.GenerateToken(
		user.GetId(),
		user.GetEmail(),
		conf.JWT.AccessTokenSecret,
		15*time.Minute,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	}

	refreshToken, err := utils.GenerateToken(
		user.GetId(),
		user.GetEmail(),
		conf.JWT.RefreshTokenSecret,
		time.Hour*24*7,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		15*60*1000,
		"/",
		"",
		conf.Server.GinMode == "release",
		true,
	)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		7*24*60*1000,
		"/",
		"",
		conf.Server.GinMode == "release",
		true,
	)

	c.JSON(http.StatusCreated, gin.H{
		"id":        user.GetId(),
		"username":  user.GetUsername(),
		"email":     user.GetEmail(),
		"createdAt": user.GetCreatedAt(),
		"updateAt":  user.GetUpdatedAt(),
	})
}
