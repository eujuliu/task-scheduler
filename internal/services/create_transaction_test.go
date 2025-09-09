package services_test

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"testing"

	"github.com/google/uuid"

	. "scheduler/test"
)

func TestCreateTransactionService(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	userData := map[string]string{
		"username": "testuser",
		"email":    "test@email.com",
		"password": "Password@123",
	}

	user, err := CreateUserService.Execute(
		userData["username"],
		userData["email"],
		userData["password"],
	)

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(
		user.GetId(),
		20,
		"BRL",
		entities.TypeTransactionPurchase,
		"",
		uuid.NewString(),
	)

	Ok(t, err)

	Equals(t, entities.StatusPending, transaction.GetStatus())
	Equals(t, entities.TypeTransactionPurchase, transaction.GetType())
}

func TestCreateTransactionService_UserNotExist(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateTransactionService.Execute(
		uuid.NewString(),
		20,
		"BRL",
		entities.TypeTransactionPurchase,
		"",
		uuid.NewString(),
	)

	Assert(t, err != nil, "expect error got success")
	Equals(t, errors.USER_NOT_FOUND_ERROR().Error(), err.Error())
}
