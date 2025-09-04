package http_handlers

import (
	"net/http"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func CancelTask(c *gin.Context) {
	taskId := c.Param("id")

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

	postgres.DB.BeginTransaction()

	task, err := updateTaskService.Cancel(taskId)
	if err != nil {
		_ = postgres.DB.RollbackTransaction()

		_ = c.Error(err)

		return
	}

	_ = postgres.DB.CommitTransaction()

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
