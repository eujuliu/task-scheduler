package http_handlers

import (
	"net/http"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateTaskRequest struct {
	RunAt    time.Time `json:"runAt"    binding:"required,date"`
	Timezone string    `json:"timezone" binding:"required,timezone"`
	Priority int       `json:"priority" binding:"required"`
}

func UpdateTask(c *gin.Context) {
	taskId := c.Param("id")
	var json UpdateTaskRequest

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

	updateTransactionService := services.NewUpdateTaskTransactionService(
		userRepository,
		transactionRepository,
		errorRepository,
	)
	updateTaskService := services.NewUpdateTaskService(
		taskRepository,
		transactionRepository,
		updateTransactionService,
	)

	task, err := updateTaskService.Execute(taskId, json.RunAt, json.Timezone, json.Priority)
	if err != nil {
		_ = c.Error(err)

		return
	}

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
