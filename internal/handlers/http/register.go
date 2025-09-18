package http_handlers

import (
	"fmt"
	"net/http"
	"scheduler/internal/config"
	"scheduler/internal/services"
	"scheduler/pkg/redis"
	"scheduler/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterHandler struct {
	config            *config.Config
	rdb               *redis.Redis
	createUserService *services.CreateUserService
}

func NewRegisterHandler(
	config *config.Config,
	rdb *redis.Redis,
	createUserService *services.CreateUserService,
) *RegisterHandler {
	return &RegisterHandler{
		config:            config,
		rdb:               rdb,
		createUserService: createUserService,
	}
}

func (h *RegisterHandler) Handle(c *gin.Context) {
	var json RegisterRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error(), "code": http.StatusBadRequest, "success": false},
		)
		return
	}

	user, err := h.createUserService.Execute(json.Username, json.Email, json.Password)
	if err != nil {
		_ = c.Error(err)

		return
	}

	accessTokenDuration := 15 * time.Minute
	refreshTokenDuration := 24 * 7 * time.Hour

	accessToken, err := utils.GenerateToken(
		user.GetId(),
		user.GetEmail(),
		h.config.JWT.AccessTokenSecret,
		accessTokenDuration,
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
		refreshTokenDuration,
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
		int(accessTokenDuration),
		"/",
		"",
		h.config.Server.GinMode == "release",
		true,
	)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(refreshTokenDuration),
		"/",
		"",
		h.config.Server.GinMode == "release",
		true,
	)

	_, err = h.rdb.Set(
		c,
		fmt.Sprintf("session_id:%v", user.GetId()),
		user.GetId(),
		accessTokenDuration,
	)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        user.GetId(),
		"username":  user.GetUsername(),
		"email":     user.GetEmail(),
		"createdAt": user.GetCreatedAt(),
		"updateAt":  user.GetUpdatedAt(),
	})
}
