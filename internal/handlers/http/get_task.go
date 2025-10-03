package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"time"

	"github.com/gin-gonic/gin"
)

type GetTaskResponse struct {
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

type GetTaskHandler struct {
	getTaskService *services.GetTaskService
}

func NewGetTaskHandler(getTaskService *services.GetTaskService) *GetTaskHandler {
	return &GetTaskHandler{
		getTaskService: getTaskService,
	}
}

// @Summary		Get task
// @Description	Get a specific task by ID
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"Task ID"
// @Success		200	{object}	GetTaskResponse
// @Failure		404	{object}	errors.Error
// @Router			/task/{id} [get]
func (h *GetTaskHandler) Handle(c *gin.Context) {
	id := c.Param("id")

	userId, _ := helpers.GetUserID(c)
	task, err := h.getTaskService.Execute(userId, id)
	if err != nil {
		_ = c.Error(err)

		return
	}

	response := GetTaskResponse{
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
	}

	c.JSON(http.StatusOK, response)
}
