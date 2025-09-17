package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/postgres"
	"time"

	"github.com/gin-gonic/gin"
)

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

func (h *UpdateTaskHandler) Handle(c *gin.Context) {
	taskId := c.Param("id")
	var json UpdateTaskRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error(), "code": http.StatusBadRequest, "success": false},
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
