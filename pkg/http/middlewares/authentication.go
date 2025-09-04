package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"scheduler/internal/config"
	"scheduler/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	cookie, err := c.Request.Cookie("access_token")
	if err != nil {
		slog.Debug("Missing access token")

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing access token",
		})

		c.Abort()
		return
	}

	if cookie.Value == "" {
		slog.Debug("Invalid access token")

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid access token",
		})

		c.Abort()
		return
	}

	claims, err := utils.ValidateToken(cookie.Value, config.Instance.JWT.AccessTokenSecret)
	if err != nil {
		slog.Debug(fmt.Sprintf("Token validation failed %s", err))

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		c.Abort()
		return
	}

	c.Set("UserID", claims.UserID)
	c.Set("Email", claims.Email)

	slog.Debug("User authenticated successfully",
		"user_id", claims.UserID,
		"email", claims.Email,
	)

	c.Next()
}

func VerifyRefreshToken(c *gin.Context) {
	cookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		slog.Debug("Missing refresh token")

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing refresh token",
		})

		c.Abort()
		return
	}

	if cookie.Value == "" {
		slog.Debug("Invalid refresh token")

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid refresh token",
		})

		c.Abort()
		return
	}

	_, err = utils.ValidateToken(cookie.Value, config.Instance.JWT.RefreshTokenSecret)
	if err != nil {
		slog.Debug(fmt.Sprintf("Token validation failed %s", err))

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		c.Abort()
		return
	}

	slog.Debug("Request token valid")

	c.Next()
}
