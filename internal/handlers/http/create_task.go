package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/postgres"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Type        string    `json:"type"        binding:"required"`
	RunAt       time.Time `json:"runAt"       binding:"required,date,utc"`
	Timezone    string    `json:"timezone"    binding:"required,timezone"`
	Priority    int       `json:"priority"    binding:"required"`
	ReferenceID string    `json:"referenceId" binding:"required,uuid4"`
}

type CreateTaskHandler struct {
	db                *postgres.Database
	createTaskService *services.CreateTaskService
}

func NewCreateTaskHandler(
	db *postgres.Database,
	createTaskService *services.CreateTaskService,
) *CreateTaskHandler {
	return &CreateTaskHandler{
		db:                db,
		createTaskService: createTaskService,
	}
}

func (h *CreateTaskHandler) Handle(c *gin.Context) {
	var json CreateTaskRequest
	idempotencyKey := c.GetHeader("Idempotency-Key")

	if idempotencyKey == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error":   "Missing idempotency key header (Idempotency-Key)",
				"code":    http.StatusBadRequest,
				"success": false,
			},
		)

		return
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error(), "code": http.StatusBadRequest, "success": false},
		)
		return
	}

	userId, ok := helpers.GetUserID(c)

	if !ok {

		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This access token is not valid",
			"success": false,
		})

		return
	}

	h.db.BeginTransaction()

	task, err := h.createTaskService.Execute(
		json.Type,
		json.RunAt,
		json.Timezone,
		json.Priority,
		userId,
		json.ReferenceID,
		idempotencyKey,
	)
	if err != nil {
		_ = h.db.RollbackTransaction()

		_ = c.Error(err)

		return
	}

	_ = h.db.CommitTransaction()

	c.JSON(http.StatusCreated, gin.H{
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
