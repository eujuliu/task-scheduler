package services_test

import (
	"scheduler/internal/entities"
	"testing"

	"github.com/google/uuid"

	. "scheduler/test"
)

func TestUpdatePurchaseTransactionService_Complete(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute(
		"testuser",
		"test@email.com",
		"Password@123",
	)

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(
		user.GetId(),
		20,
		"USD",
		entities.TypeTransactionPurchase,
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	_, err = UpdatePurchaseTransactionService.Complete(transaction.GetId())

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 20, user.GetCredits())
}

func TestUpdatePurchaseTransactionService_Fail(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute(
		"testuser",
		"test@email.com",
		"Password@123",
	)

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(
		user.GetId(),
		20,
		"USD",
		entities.TypeTransactionPurchase,
		uuid.NewString(),
		uuid.NewString(),
	)

	Ok(t, err)

	transaction, err = UpdatePurchaseTransactionService.Fail(
		transaction.GetId(),
		"Credit card invalid",
	)

	Ok(t, err)
	Equals(t, entities.StatusFailed, transaction.GetStatus())

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())
}
