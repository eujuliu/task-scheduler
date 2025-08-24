package services_test

import (
	"fmt"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	. "scheduler/test"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateTaskService(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	task, err := CreateTaskService.Execute("email", time.Now().AddDate(0, 1, 0), "America/Sao_Paulo", entities.PriorityMedium, user.GetId(), uuid.NewString(), uuid.NewString())

	Ok(t, err)

	transaction, err := TransactionRepository.GetFirstByReferenceId(fmt.Sprintf("task_%s", task.GetReferenceId()))

	Ok(t, err)

	Equals(t, fmt.Sprintf("task_%s", task.GetReferenceId()), transaction.GetReferenceId())
	Equals(t, task.GetIdempotencyKey(), transaction.GetIdempotencyKey())
	Equals(t, "America/Sao_Paulo", task.GetTimezone())
	Equals(t, entities.PriorityMedium, task.GetPriority())
}

func TestCreateTaskService_UserNotFound(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateTaskService.Execute("email", time.Now().AddDate(0, 1, 0), "America/Sao_Paulo", entities.PriorityMedium, uuid.NewString(), uuid.NewString(), uuid.NewString())

	Assert(t, err != nil, "expected error got success")
	Equals(t, errors.USER_NOT_FOUND_ERROR().Error(), err.Error())
}
