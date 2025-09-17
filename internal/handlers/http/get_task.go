package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"

	"github.com/gin-gonic/gin"
)

type GetTaskHandler struct {
	getTaskService *services.GetTaskService
}

func NewGetTaskHandler(getTaskService *services.GetTaskService) *GetTaskHandler {
	return &GetTaskHandler{
		getTaskService: getTaskService,
	}
}

func (h *GetTaskHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	userId, ok := helpers.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This access token is not valid",
			"success": false,
		})
	}

	task, err := h.getTaskService.Execute(userId, id)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          task.GetId(),
		"status":      task.GetStatus(),
		"cost":        task.GetCost(),
		"runAt":       task.GetRunAt(),
		"timezone":    task.GetTimezone(),
		"retries":     task.GetRetries(),
		"priority":    task.GetPriority(),
		"type":        task.GetType(),
		"referenceId": task.GetReferenceId(),
		"createdAt":   task.GetCreatedAt(),
		"updateAt":    task.GetUpdatedAt(),
	})
}
