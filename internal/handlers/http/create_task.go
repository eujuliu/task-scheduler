package http_handlers

import (
	"net/http"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/postgres"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Type        string    `json:"type"        binding:"required"`
	RunAt       time.Time `json:"runAt"       binding:"required,date"`
	Timezone    string    `json:"timezone"    binding:"required,timezone"`
	Priority    int       `json:"priority"    binding:"required"`
	ReferenceID string    `json:"referenceId" binding:"required,uuid4"`
}

func CreateTask(c *gin.Context) {
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

	userRepository := postgres_repos.NewPostgresUserRepository()
	transactionRepository := postgres_repos.NewPostgresTransactionRepository()
	taskRepository := postgres_repos.NewPostgresTaskRepository()
	errorRepository := postgres_repos.NewPostgresErrorRepository()
	createTransactionService := services.NewCreateTransactionService(
		userRepository,
		transactionRepository,
	)
	updateTransactionService := services.NewUpdateTaskTransactionService(
		userRepository,
		transactionRepository,
		errorRepository,
	)

	createTaskService := services.NewCreateTaskService(
		userRepository,
		taskRepository,
		createTransactionService,
		updateTransactionService,
	)

	userId, ok := helpers.GetUserID(c)

	if !ok {

		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This access token is not valid",
			"success": false,
		})

		return
	}

	postgres.DB.BeginTransaction()

	task, err := createTaskService.Execute(
		json.Type,
		json.RunAt,
		json.Timezone,
		json.Priority,
		userId,
		json.ReferenceID,
		idempotencyKey,
	)
	if err != nil {
		_ = postgres.DB.RollbackTransaction()

		_ = c.Error(err)

		return
	}

	_ = postgres.DB.CommitTransaction()

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
