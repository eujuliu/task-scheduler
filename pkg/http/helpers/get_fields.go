package helpers

import (
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("UserID")
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("Email")
	if !exists {
		return "", false
	}

	mail, ok := email.(string)
	return mail, ok
}
