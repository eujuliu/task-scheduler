package services_test

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"testing"
	"time"

	"github.com/google/uuid"

	. "scheduler/test"
)

func TestCreateTaskService(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute(
		"testuser",
		"test@email.com",
		"Password@123",
	)

	Ok(t, err)

	user.AddCredits(10)

	err = UserRepository.Update(user)

	Ok(t, err)

	task, err := CreateTaskService.Execute(
		"email",
		time.Now().AddDate(0, 1, 0),
		"America/Sao_Paulo",
		entities.PriorityMedium,
		user.GetId(),
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	transaction, err := TransactionRepository.GetFirstByReferenceId(task.GetReferenceId())

	Ok(t, err)

	user, err = UserRepository.GetFirstById(user.GetId())

	Ok(t, err)

	Equals(t, "America/Sao_Paulo", task.GetTimezone())
	Equals(t, entities.PriorityMedium, task.GetPriority())
	Equals(t, task.GetCost(), transaction.GetCredits())
	Equals(t, task.GetIdempotencyKey(), transaction.GetIdempotencyKey())
	Equals(t, 0, user.GetCredits())
	Equals(t, 10, user.GetFrozenCredits())
}

func TestCreateTaskService_UserNotFound(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateTaskService.Execute(
		"email",
		time.Now().AddDate(0, 1, 0),
		"America/Sao_Paulo",
		entities.PriorityMedium,
		uuid.NewString(),
		uuid.NewString(),
		uuid.NewString(),
	)

	Assert(t, err != nil, "expected error got success")
	Equals(t, errors.USER_NOT_FOUND_ERROR().Error(), err.Error())
}
