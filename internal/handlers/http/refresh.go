package http_handlers

import (
	"net/http"
	"scheduler/internal/config"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func Refresh(c *gin.Context) {
	userId, ok := helpers.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid claim into token",
		})
	}

	email, ok := helpers.GetEmail(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid claim into token",
		})
	}

	accessToken, err := utils.GenerateToken(
		userId,
		email,
		config.Instance.JWT.AccessTokenSecret,
		15*time.Minute,
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
		config.Instance.Server.GinMode == "release",
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "New access Token added to Cookies",
	})
}
