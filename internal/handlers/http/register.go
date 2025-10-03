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

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
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

// @Summary		Register a new user
// @Description	Register a new user with username, email, and password
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		RegisterRequest	true	"Register request"
// @Success		201		{object}	RegisterResponse
// @Failure		400		{object}	errors.Error
// @Failure		404		{object}	errors.Error
// @Router			/auth/register [post]
func (h *RegisterHandler) Handle(c *gin.Context) {
	var json RegisterRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"id": uuid.NewString(), "error": err.Error(), "code": http.StatusBadRequest},
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

	response := RegisterResponse{
		User: UserResponse{
			ID:        user.GetId(),
			Username:  user.GetUsername(),
			Email:     user.GetEmail(),
			CreatedAt: user.GetCreatedAt().Format(time.RFC3339),
			UpdatedAt: user.GetUpdatedAt().Format(time.RFC3339),
		},
		Token: accessToken,
	}

	c.JSON(http.StatusCreated, response)
}
