package services_test

import (
	"scheduler/internal/entities"
	. "scheduler/test"
	"strconv"
	"testing"

	"github.com/google/uuid"
)

func TestGetTransactionsService(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	user.AddCredits(200)

	err = UserRepository.Update(user)

	Ok(t, err)

	for i := range 10 {
		transaction, err := entities.NewTransaction(
			user.GetId(),
			10,
			5,
			"BRL",
			entities.TypeTransactionPurchase,
			uuid.NewString(),
			strconv.Itoa(i+1),
		)

		Ok(t, err)

		_ = TransactionRepository.Create(transaction)
	}

	transactions := GetTransactionsService.Execute(user.GetId())

	Equals(t, 10, len(transactions))
	Equals(t, 10, transactions[4].GetCredits())
	Equals(t, "10", transactions[9].GetIdempotencyKey())
}
