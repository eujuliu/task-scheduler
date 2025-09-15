package services_test

import (
	"scheduler/internal/entities"
	. "scheduler/test"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetTaskService(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	user.AddCredits(100)

	err = UserRepository.Update(user)

	Ok(t, err)

	taskOld, err := CreateTaskService.Execute(
		"video",
		time.Now().Add(10*time.Minute),
		"America/Sao_Paulo",
		entities.PriorityMedium,
		user.GetId(),
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	task, err := GetTaskService.Execute(user.GetId(), taskOld.GetId())

	Ok(t, err)
	Equals(t, "video", task.GetType())
}
