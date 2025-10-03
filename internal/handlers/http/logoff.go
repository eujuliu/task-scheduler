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

// @Summary		Log off user
// @Description	Log off the user by invalidating session
// @Tags			auth
// @Accept			json
// @Produce		json
// @Success		200
// @Router			/logoff [delete]
func (h *LogoffHandler) Handle(c *gin.Context) {
	userId, _ := helpers.GetUserID(c)

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
