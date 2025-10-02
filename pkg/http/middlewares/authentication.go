package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"scheduler/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	header := c.GetHeader("Authorization")

	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing authorization header",
		})

		c.Abort()
		return
	}

	if !strings.HasPrefix(header, "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authorization token",
		})

		c.Abort()
		return
	}

	token := strings.Split(header, " ")[1]

	claims, err := utils.ValidateToken(token, utils.GetEnv("ACCESS_TOKEN_SECRET", ""))
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

	_, err = utils.ValidateToken(cookie.Value, utils.GetEnv("REFRESH_TOKEN_SECRET", ""))
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
