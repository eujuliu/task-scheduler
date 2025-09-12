package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/postgres"

	"github.com/gin-gonic/gin"
)

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
