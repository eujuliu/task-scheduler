package services_test

import (
	"scheduler/internal/errors"
	. "scheduler/test"
	"testing"

	"github.com/google/uuid"
)

func TestUpdateTransactionService_Purchase(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(user.GetId(), 20, 10, "USD", "purchase", uuid.NewString(), uuid.NewString())

	Ok(t, err)

	_, err = UpdateTransactionService.Execute(transaction.GetId(), "completed", make(map[string]any))

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 20, user.GetCredits())
}

func TestUpdateTransactionService_Task(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)
	Equals(t, 0, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())

	user.AddCredits(20)
	err = user.AddFrozenCredits(20)

	Ok(t, err)

	err = UserRepository.Update(user)

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(user.GetId(), 20, 0, "", "task_send", uuid.NewString(), uuid.NewString())

	Ok(t, err)

	_, err = UpdateTransactionService.Execute(transaction.GetId(), "completed", make(map[string]any))

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())
}

func TestUpdateTransactionService_PurchaseFailed(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)
	Equals(t, 0, user.GetCredits())

	transaction, err := CreateTransactionService.Execute(user.GetId(), 20, 10, "USD", "purchase", uuid.NewString(), uuid.NewString())

	Ok(t, err)

	var options = map[string]any{
		"reason": "The payment getaway don't finish the operation correctly",
	}

	_, err = UpdateTransactionService.Execute(transaction.GetId(), "failed", options)

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)
	Equals(t, 0, user.GetCredits())
}

func TestUpdateTransactionService_TaskFailedRefund(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	user.AddCredits(20)

	err = UserRepository.Update(user)

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(user.GetId(), 20, 0, "", "task_send", uuid.NewString(), uuid.NewString())

	Ok(t, err)

	transaction, err = UpdateTransactionService.Execute(transaction.GetId(), "frozen", make(map[string]any))

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())
	Equals(t, 20, user.GetFrozenCredits())

	var options = map[string]any{
		"refund": true,
		"reason": "Email provider not working!",
	}

	_, err = UpdateTransactionService.Execute(transaction.GetId(), "failed", options)

	Ok(t, err)

	user, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)

	Equals(t, 20, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())
}

func TestUpdateTransactionService_TransactionNotFound(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := UpdateTransactionService.Execute(uuid.NewString(), "completed", make(map[string]any))

	Assert(t, err != nil, "expected error got success")
	Equals(t, errors.TRANSACTION_NOT_FOUND().Error(), err.Error())
}

func TestUpdateTransactionService_UserNotExist(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	transaction, err := CreateTransactionService.Execute(user.GetId(), 20, 10, "USD", "purchase", uuid.NewString(), uuid.NewString())

	Ok(t, err)

	err = UserRepository.Delete(user.GetId())

	Ok(t, err)

	_, err = UpdateTransactionService.Execute(transaction.GetId(), "completed", make(map[string]any))

	Assert(t, err != nil, "expected error got success")
	Equals(t, errors.USER_NOT_FOUND_ERROR().Error(), err.Error())
}
