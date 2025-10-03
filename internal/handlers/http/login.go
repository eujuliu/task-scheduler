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
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Credits       int    `json:"credits,omitempty"`
	FrozenCredits int    `json:"frozen_credits,omitempty"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type LoginHandler struct {
	config         *config.Config
	rdb            *redis.Redis
	getUserService *services.GetUserService
}

func NewLoginHandler(
	config *config.Config,
	rdb *redis.Redis,
	getUserService *services.GetUserService,
) *LoginHandler {
	return &LoginHandler{
		config:         config,
		rdb:            rdb,
		getUserService: getUserService,
	}
}

// @Summary		Login user
// @Description	Authenticate user with email and password
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		LoginRequest	true	"Login request"
// @Success		200		{object}	LoginResponse
// @Failure		400		{object}	errors.Error
// @Failure		404		{object}	errors.Error
// @Router			/auth/login [post]
func (h *LoginHandler) Handle(c *gin.Context) {
	var json LoginRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"id": uuid.NewString(), "error": err.Error(), "code": http.StatusBadRequest},
		)
		return
	}

	user, err := h.getUserService.Execute(json.Email, json.Password)
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
			"id":    uuid.NewString(),
			"code":  http.StatusUnauthorized,
			"error": err.Error(),
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
			"id":    uuid.NewString(),
			"code":  http.StatusUnauthorized,
			"error": err.Error(),
		})

		return
	}

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

	response := LoginResponse{
		User: UserResponse{
			ID:            user.GetId(),
			Username:      user.GetUsername(),
			Email:         user.GetEmail(),
			Credits:       user.GetCredits(),
			FrozenCredits: user.GetFrozenCredits(),
			CreatedAt:     user.GetCreatedAt().Format(time.RFC3339),
			UpdatedAt:     user.GetUpdatedAt().Format(time.RFC3339),
		},
		Token: accessToken,
	}

	c.JSON(http.StatusOK, response)
}
