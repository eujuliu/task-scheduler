package http_handlers

import (
	"fmt"
	"net/http"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/redis"

	"github.com/gin-gonic/gin"
)

type LogoffHandler struct {
	rdb *redis.Redis
}

func NewLogoffHandler(rdb *redis.Redis) *LogoffHandler {
	return &LogoffHandler{
		rdb: rdb,
	}
}

func (h *LogoffHandler) Handle(c *gin.Context) {
	userId, ok := helpers.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"error":   "Invalid claim into token",
			"success": false,
		})
	}

	_, err := h.rdb.Del(
		c,
		fmt.Sprintf("session_id:%v", userId),
	)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.Status(http.StatusOK)
}
