package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetTasksHandler struct {
	getTasksService *services.GetTasksByUserIdService
}

func NewGetTasksHandler(getTasksService *services.GetTasksByUserIdService) *GetTasksHandler {
	return &GetTasksHandler{
		getTasksService: getTasksService,
	}
}

func (h *GetTasksHandler) Handle(c *gin.Context) {
	userId, ok := helpers.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This access token is not valid",
			"success": false,
		})
	}

	var offset int
	var limit int
	orderBy, _ := c.GetQuery("orderBy")

	if value, ok := c.GetQuery("offset"); ok {
		v, _ := strconv.Atoi(value)

		offset = v
	}

	if value, ok := c.GetQuery("limit"); ok {
		v, _ := strconv.Atoi(value)

		limit = v
	}

	tasks := h.getTasksService.Execute(userId, &offset, &limit, &orderBy)
	result := []map[string]any{}

	for _, task := range tasks {
		result = append(result, map[string]any{
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

	c.JSON(http.StatusOK, result)
}
