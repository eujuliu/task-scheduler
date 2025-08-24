package services_test

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	. "scheduler/test"
	"testing"

	"github.com/google/uuid"
)

func TestCreateTransactionService(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	var userData = map[string]string{
		"username": "testuser",
		"email":    "test@email.com",
		"password": "Password@123",
	}

	user, err := CreateUserService.Execute(userData["username"], userData["email"], userData["password"])

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(user.GetId(), 20, 10, "BRL", entities.TypeTransactionPurchase, uuid.NewString(), uuid.NewString())

	Ok(t, err)

	Equals(t, entities.StatusPending, transaction.GetStatus())
	Equals(t, entities.TypeTransactionPurchase, transaction.GetType())
}

func TestCreateTransactionService_UserNotExist(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateTransactionService.Execute(uuid.NewString(), 20, 10, "BRL", entities.TypeTransactionPurchase, uuid.NewString(), uuid.NewString())

	Assert(t, err != nil, "expect error got success")
	Equals(t, errors.USER_NOT_FOUND_ERROR().Error(), err.Error())
}

func TestCreateTransactionService_DuplicatedReference(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	var userData = map[string]string{
		"username": "testuser",
		"email":    "test@email.com",
		"password": "Password@123",
	}

	var referenceId = uuid.NewString()
	userOld, err := CreateUserService.Execute(userData["username"], userData["email"], userData["password"])

	Ok(t, err)

	_, err = CreateTransactionService.Execute(userOld.GetId(), 20, 10, "BRL", entities.TypeTransactionPurchase, referenceId, uuid.NewString())

	Ok(t, err)

	_, err = CreateTransactionService.Execute(userOld.GetId(), 20, 10, "BRL", entities.TypeTransactionPurchase, referenceId, uuid.NewString())

	Assert(t, err != nil, "expect error got success")
	Equals(t, errors.TRANSACTION_ALREADY_EXISTS_ERROR().Error(), err.Error())
}

func TestCreateTransactionService_DuplicatedIdempotency(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	var userData = map[string]string{
		"username": "testuser",
		"email":    "test@email.com",
		"password": "Password@123",
	}

	var idempotencyKey = uuid.NewString()
	userOld, err := CreateUserService.Execute(userData["username"], userData["email"], userData["password"])

	Ok(t, err)

	_, err = CreateTransactionService.Execute(userOld.GetId(), 20, 10, "BRL", entities.TypeTransactionPurchase, uuid.NewString(), idempotencyKey)

	Ok(t, err)

	_, err = CreateTransactionService.Execute(userOld.GetId(), 20, 10, "BRL", entities.TypeTransactionPurchase, uuid.NewString(), idempotencyKey)

	Assert(t, err != nil, "expect error got success")
	Equals(t, errors.TRANSACTION_ALREADY_EXISTS_ERROR().Error(), err.Error())
}
