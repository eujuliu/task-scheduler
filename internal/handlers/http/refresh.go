package http_handlers

import (
	"net/http"
	"scheduler/internal/config"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type RefreshTokenHandler struct {
	config *config.Config
}

func NewRefreshTokenHandler(config *config.Config) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		config: config,
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

	c.SetCookie(
		"access_token",
		accessToken,
		15*60*1000,
		"/",
		"",
		h.config.Server.GinMode == "release",
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "New access Token added to Cookies",
	})
}
