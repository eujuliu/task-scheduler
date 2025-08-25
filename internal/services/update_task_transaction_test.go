package services_test

import (
	"scheduler/internal/entities"
	"testing"

	"github.com/google/uuid"

	. "scheduler/test"
)

func TestUpdateTaskTransactionService_Complete(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute(
		"testuser",
		"test@email.com",
		"Password@123",
	)

	Ok(t, err)
	Equals(t, 0, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())

	user.AddCredits(20)
	err = user.AddFrozenCredits(20)

	Ok(t, err)

	err = UserRepository.Update(user)

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(
		user.GetId(),
		20,
		0,
		"",
		entities.TypeTransactionTaskSend,
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	_, err = UpdateTaskTransactionService.Complete(transaction.GetId())

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())
}

func TestUpdateTaskTransactionService_Frozen(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute(
		"testuser",
		"test@email.com",
		"Password@123",
	)

	Ok(t, err)

	user.AddCredits(20)

	err = UserRepository.Update(user)

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(
		user.GetId(),
		20,
		0,
		"",
		entities.TypeTransactionTaskSend,
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	transaction, err = UpdateTaskTransactionService.Frozen(transaction.GetId())

	Ok(t, err)
	Equals(t, entities.StatusFrozen, transaction.GetStatus())

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)
	Equals(t, 0, user.GetCredits())
	Equals(t, 20, user.GetFrozenCredits())
}

func TestUpdateTaskTransactionService_FailWithRefund(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute(
		"testuser",
		"test@email.com",
		"Password@123",
	)

	Ok(t, err)

	user.AddCredits(20)

	err = UserRepository.Update(user)

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(
		user.GetId(),
		20,
		0,
		"",
		entities.TypeTransactionTaskSend,
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	transaction, err = UpdateTaskTransactionService.Frozen(transaction.GetId())

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())
	Equals(t, 20, user.GetFrozenCredits())

	_, err = UpdateTaskTransactionService.Fail(
		transaction.GetId(),
		true,
		"Email service don't working",
	)

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 20, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())
}

func TestUpdateTaskTransactionService_FailWithoutRefund(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute(
		"testuser",
		"test@email.com",
		"Password@123",
	)

	Ok(t, err)

	user.AddCredits(20)

	err = UserRepository.Update(user)

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(
		user.GetId(),
		20,
		0,
		"",
		entities.TypeTransactionTaskSend,
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	transaction, err = UpdateTaskTransactionService.Frozen(transaction.GetId())

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())
	Equals(t, 20, user.GetFrozenCredits())

	_, err = UpdateTaskTransactionService.Fail(
		transaction.GetId(),
		false,
		"Email service don't working",
	)

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())
}

func TestUpdateTaskTransactionService_Cancel(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute(
		"testuser",
		"test@email.com",
		"Password@123",
	)

	Ok(t, err)

	user.AddCredits(20)

	err = UserRepository.Update(user)

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(
		user.GetId(),
		20,
		0,
		"",
		entities.TypeTransactionTaskSend,
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	transaction, err = UpdateTaskTransactionService.Frozen(transaction.GetId())

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())
	Equals(t, 20, user.GetFrozenCredits())

	_, err = UpdateTaskTransactionService.Cancel(transaction.GetId())

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 20, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())
}
