package services_test

import (
	"scheduler/internal/entities"
	. "scheduler/test"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetTasksByRunAtService_Asc(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	runAt := time.Now()

	for i := range 10 {
		task, _ := entities.NewTask(
			"email",
			user.GetId(),
			10,
			runAt.AddDate(0, 3, i+1),
			"America/Sao_Paulo",
			uuid.NewString(),
			strconv.Itoa(i),
		)

		_ = TaskRepository.Create(task)
	}

	tasks := GetTasksByRunAtService.Execute(entities.StatusPending, true, 10, nil)

	Ok(t, err)
	Equals(t, 10, len(tasks))
	Equals(t, runAt.AddDate(0, 3, 1), tasks[0].GetRunAt())
}
