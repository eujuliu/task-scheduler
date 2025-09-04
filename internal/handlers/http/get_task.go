package http_handlers

import (
	"net/http"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/http/helpers"

	"github.com/gin-gonic/gin"
)

func GetTask(c *gin.Context) {
	id := c.Param("id")

	userRepository := postgres_repos.NewPostgresUserRepository()
	taskRepository := postgres_repos.NewPostgresTaskRepository()
	getTaskService := services.NewGetTaskService(
		userRepository,
		taskRepository,
	)

	userId, ok := helpers.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This access token is not valid",
			"success": false,
		})
	}

	task, err := getTaskService.Execute(userId, id)
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
