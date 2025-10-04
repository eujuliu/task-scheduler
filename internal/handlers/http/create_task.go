package http_handlers

import (
	"net/http"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/postgres"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	Type        string    `json:"type"        binding:"required"`
	RunAt       time.Time `json:"runAt"       binding:"required,date,utc"`
	Timezone    string    `json:"timezone"    binding:"required,timezone"`
	Priority    int       `json:"priority"    binding:"required"`
	ReferenceID string    `json:"referenceId" binding:"required,uuid4"`
}

type CreateTaskResponse struct {
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

// @Summary		Create task
// @Description	Create a new task for the user
// @Tags			tasks
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			Idempotency-Key	header		string				true	"Idempotency Key"
// @Param			request			body		CreateTaskRequest	true	"Create task request"
// @Success		201				{object}	CreateTaskResponse
// @Failure		400				{object}	errors.Error
// @Failure		404				{object}	errors.Error
// @Router			/task [post]
func (h *CreateTaskHandler) Handle(c *gin.Context) {
	var json CreateTaskRequest
	idempotencyKey := c.GetHeader("Idempotency-Key")

	if idempotencyKey == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"id":    uuid.NewString(),
				"error": "Missing idempotency key header (Idempotency-Key)",
				"code":  http.StatusBadRequest,
			},
		)

		return
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"id": uuid.NewString(), "error": err.Error(), "code": http.StatusBadRequest},
		)
		return
	}

	userId, _ := helpers.GetUserID(c)

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

	response := CreateTaskResponse{
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

	c.JSON(http.StatusCreated, response)
}
