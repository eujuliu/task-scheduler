package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Cost        int    `json:"cost"`
	RunAt       string `json:"runAt"`
	Timezone    string `json:"timezone"`
	Retries     int    `json:"retries"`
	Priority    int    `json:"priority"`
	Type        string `json:"type"`
	ReferenceID string `json:"referenceId"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type GetTasksHandler struct {
	getTasksService *services.GetTasksByUserIdService
}

func NewGetTasksHandler(getTasksService *services.GetTasksByUserIdService) *GetTasksHandler {
	return &GetTasksHandler{
		getTasksService: getTasksService,
	}
}

// @Summary		Get tasks
// @Description	Get list of tasks for the user
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			offset	query		int		false	"Offset"
// @Param			limit	query		int		false	"Limit"
// @Param			orderBy	query		string	false	"Order by"
// @Success		200		{array}		TaskResponse
// @Failure		404		{object}	errors.Error
// @Router			/tasks [get]
func (h *GetTasksHandler) Handle(c *gin.Context) {
	userId, _ := helpers.GetUserID(c)

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
	result := []TaskResponse{}

	for _, task := range tasks {
		result = append(result, TaskResponse{
			ID:          task.GetId(),
			Status:      task.GetStatus(),
			Cost:        task.GetCost(),
			RunAt:       task.GetRunAt().Format(time.RFC3339),
			Timezone:    task.GetTimezone(),
			Retries:     task.GetRetries(),
			Priority:    task.GetPriority(),
			Type:        task.GetType(),
			ReferenceID: task.GetReferenceId(),
			CreatedAt:   task.GetCreatedAt().Format(time.RFC3339),
			UpdatedAt:   task.GetUpdatedAt().Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, result)
}
