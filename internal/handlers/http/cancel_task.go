package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/postgres"
	"time"

	"github.com/gin-gonic/gin"
)

type CancelTaskResponse struct {
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

type CancelTaskHandler struct {
	db                *postgres.Database
	updateTaskService *services.UpdateTaskService
}

func NewCancelTaskHandler(
	db *postgres.Database,
	updateTaskService *services.UpdateTaskService,
) *CancelTaskHandler {
	return &CancelTaskHandler{
		db:                db,
		updateTaskService: updateTaskService,
	}
}

// @Summary		Cancel task
// @Description	Cancel an existing task
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"Task ID"
// @Success		200	{object}	CancelTaskResponse
// @Failure		404	{object}	errors.Error
// @Router			/task/cancel/{id} [put]
func (h *CancelTaskHandler) Handle(c *gin.Context) {
	taskId := c.Param("id")

	h.db.BeginTransaction()

	task, err := h.updateTaskService.Cancel(taskId)
	if err != nil {
		_ = h.db.RollbackTransaction()

		_ = c.Error(err)

		return
	}

	_ = h.db.CommitTransaction()

	response := CancelTaskResponse{
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
