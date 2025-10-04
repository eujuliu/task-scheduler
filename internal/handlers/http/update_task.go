package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/postgres"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateTaskResponse struct {
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

type UpdateTaskRequest struct {
	RunAt    time.Time `json:"runAt"    binding:"required,date,utc"`
	Timezone string    `json:"timezone" binding:"required,timezone"`
	Priority int       `json:"priority" binding:"required"`
}

type UpdateTaskHandler struct {
	db                *postgres.Database
	updateTaskService *services.UpdateTaskService
}

func NewUpdateTaskHandler(
	db *postgres.Database,
	updateTaskService *services.UpdateTaskService,
) *UpdateTaskHandler {
	return &UpdateTaskHandler{
		db:                db,
		updateTaskService: updateTaskService,
	}
}

// @Summary		Update task
// @Description	Update an existing task's runAt, timezone, and priority
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Param			id		path		string				true	"Task ID"
// @Param			request	body		UpdateTaskRequest	true	"Update task request"
// @Success		200		{object}	UpdateTaskResponse
// @Failure		400		{object}	errors.Error
// @Failure		404		{object}	errors.Error
// @Router			/task/{id} [put]
func (h *UpdateTaskHandler) Handle(c *gin.Context) {
	taskId := c.Param("id")
	var json UpdateTaskRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"id": uuid.NewString(), "error": err.Error(), "code": http.StatusBadRequest},
		)
		return
	}

	h.db.BeginTransaction()

	task, err := h.updateTaskService.Execute(
		taskId,
		nil,
		&json.RunAt,
		&json.Timezone,
		&json.Priority,
	)
	if err != nil {
		_ = h.db.RollbackTransaction()
		_ = c.Error(err)

		return
	}

	_ = h.db.CommitTransaction()

	response := UpdateTaskResponse{
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
