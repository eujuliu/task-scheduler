package http_handlers

import (
	"net/http"
	"scheduler/internal/config"
	"scheduler/internal/services"
	"scheduler/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginHandler struct {
	config         *config.Config
	getUserService *services.GetUserService
}

func NewLoginHandler(config *config.Config, getUserService *services.GetUserService) *LoginHandler {
	return &LoginHandler{
		config:         config,
		getUserService: getUserService,
	}
}

func (h *LoginHandler) Handle(c *gin.Context) {
	var json LoginRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error(), "code": http.StatusBadRequest, "success": false},
		)
		return
	}

	user, err := h.getUserService.Execute(json.Email, json.Password)
	if err != nil {
		_ = c.Error(err)

		return
	}

	accessToken, err := utils.GenerateToken(
		user.GetId(),
		user.GetEmail(),
		h.config.JWT.AccessTokenSecret,
		15*time.Minute,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"error":   err.Error(),
			"success": false,
		})

		return
	}

	refreshToken, err := utils.GenerateToken(
		user.GetId(),
		user.GetEmail(),
		h.config.JWT.RefreshTokenSecret,
		time.Hour*24*7,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"error":   err.Error(),
			"success": false,
		})

		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		15*60*1000,
		"/",
		"",
		h.config.Server.GinMode == "release",
		true,
	)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		7*24*60*1000,
		"/",
		"",
		h.config.Server.GinMode == "release",
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"id":             user.GetId(),
		"username":       user.GetUsername(),
		"email":          user.GetEmail(),
		"credits":        user.GetCredits(),
		"frozen_credits": user.GetFrozenCredits(),
		"createdAt":      user.GetCreatedAt(),
		"updateAt":       user.GetUpdatedAt(),
	})
}
