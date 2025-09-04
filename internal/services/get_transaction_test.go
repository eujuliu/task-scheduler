package services_test

import (
	"scheduler/internal/entities"
	. "scheduler/test"
	"testing"

	"github.com/google/uuid"
)

func TestGetTransactionService(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	user.AddCredits(10)

	err = UserRepository.Update(user)

	Ok(t, err)

	transaction, err := entities.NewTransaction(
		user.GetId(),
		10,
		5,
		"BRL",
		entities.TypeTransactionPurchase,
		uuid.NewString(),
		"123",
	)

	Ok(t, err)

	err = TransactionRepository.Create(transaction)

	Ok(t, err)

	transaction, err = GetTransactionService.Execute(user.GetId(), transaction.GetId())

	Ok(t, err)
	Equals(t, 10, transaction.GetCredits())
	Equals(t, entities.TypeTransactionPurchase, transaction.GetType())
	Equals(t, "123", transaction.GetIdempotencyKey())
}
