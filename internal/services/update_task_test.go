package services_test

import (
	"scheduler/internal/entities"
	. "scheduler/test"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUpdateTaskService_Properties(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	user.AddCredits(20)

	err = UserRepository.Update(user)

	Ok(t, err)

	task, err := CreateTaskService.Execute(
		"email",
		time.Now().Add(10*time.Minute),
		"America/Sao_Paulo",
		entities.PriorityMedium,
		user.GetId(),
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	runAt := time.Now().Add(10 * time.Minute)

	timezone := "Europe/Berlin"
	priority := entities.PriorityHigh

	task, err = UpdateTaskService.Execute(
		task.GetId(),
		nil,
		&runAt,
		&timezone,
		&priority,
	)

	Ok(t, err)
	Equals(t, "Europe/Berlin", task.GetTimezone())
	Equals(t, runAt, task.GetRunAt())
}

func TestUpdateTaskService_Complete(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	user.AddCredits(20)

	err = UserRepository.Update(user)

	Ok(t, err)

	task, err := CreateTaskService.Execute(
		"email",
		time.Now().Add(10*time.Minute),
		"America/Sao_Paulo",
		entities.PriorityMedium,
		user.GetId(),
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	task, err = UpdateTaskService.Complete(task.GetId())

	Ok(t, err)
	Equals(t, entities.StatusCompleted, task.GetStatus())
}

func TestUpdateTaskService_FailWithRefund(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	user.AddCredits(20)

	err = UserRepository.Update(user)

	Ok(t, err)

	task, err := CreateTaskService.Execute(
		"email",
		time.Now().Add(10*time.Minute),
		"America/Sao_Paulo",
		entities.PriorityMedium,
		user.GetId(),
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	task, err = UpdateTaskService.Fail(task.GetId(), true, "Email provider not working")

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)
	Equals(t, 20, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())
	Equals(t, entities.StatusFailed, task.GetStatus())
}

func TestUpdateTaskService_FailWithoutRefund(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	user.AddCredits(10)

	err = UserRepository.Update(user)

	Ok(t, err)

	task, err := CreateTaskService.Execute(
		"email",
		time.Now().Add(10*time.Minute),
		"America/Sao_Paulo",
		entities.PriorityMedium,
		user.GetId(),
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	task, err = UpdateTaskService.Fail(task.GetId(), false, "Audience email with wrong format")

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)
	Equals(t, 0, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())
	Equals(t, entities.StatusFailed, task.GetStatus())
}

func TestUpdateTaskService_Cancel(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	user.AddCredits(10)

	err = UserRepository.Update(user)

	Ok(t, err)

	task, err := CreateTaskService.Execute(
		"email",
		time.Now().Add(10*time.Minute),
		"America/Sao_Paulo",
		entities.PriorityMedium,
		user.GetId(),
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	task, err = UpdateTaskService.Cancel(task.GetId())

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)
	Equals(t, 10, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())
	Equals(t, entities.StatusCanceled, task.GetStatus())
}
