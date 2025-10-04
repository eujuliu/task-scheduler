package http_handlers

import (
	"fmt"
	"net/http"
	"scheduler/internal/config"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/redis"
	"scheduler/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

type RefreshTokenHandler struct {
	config *config.Config
	rdb    *redis.Redis
}

func NewRefreshTokenHandler(config *config.Config, rdb *redis.Redis) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		config: config,
		rdb:    rdb,
	}
}

// @Summary		Refresh access token
// @Description	Refresh the access token using refresh token
// @Tags			auth
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Security		RefreshTokenAuth
// @Success		200	{object}	RefreshTokenResponse
// @Failure		401	{object}	errors.Error
// @Router			/refresh [post]
func (h *RefreshTokenHandler) Handle(c *gin.Context) {
	userId, _ := helpers.GetUserID(c)
	email, _ := helpers.GetEmail(c)

	accessTokenDuration := 15 * time.Minute

	accessToken, err := utils.GenerateToken(
		userId,
		email,
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

	_, err = h.rdb.Set(
		c,
		fmt.Sprintf("session_id:%v", userId),
		userId,
		accessTokenDuration,
	)
	if err != nil {
		_ = c.Error(err)

		return
	}

	response := RefreshTokenResponse{
		Token: accessToken,
	}

	c.JSON(http.StatusOK, response)
}
