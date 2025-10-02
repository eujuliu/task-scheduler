package http_handlers

import (
	"net/http"
	"scheduler/internal/config"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/redis"
	"scheduler/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

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

func (h *RefreshTokenHandler) Handle(c *gin.Context) {
	userId, ok := helpers.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"error":   "Invalid claim into token",
			"success": false,
		})
	}

	email, ok := helpers.GetEmail(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"error":   "Invalid claim into token",
			"success": false,
		})
	}

	accessToken, err := utils.GenerateToken(
		userId,
		email,
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

	c.JSON(http.StatusOK, gin.H{
		"token": accessToken,
	})
	c.Status(http.StatusOK)
}
