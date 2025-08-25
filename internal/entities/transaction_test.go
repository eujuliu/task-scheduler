package entities_test

import (
	"scheduler/internal/errors"
	"testing"

	. "scheduler/internal/entities"

	"github.com/google/uuid"

	. "scheduler/test"
)

func TestTransaction_New(t *testing.T) {
	transaction, err := NewTransaction(
		uuid.NewString(),
		40,
		20,
		"BRL",
		TypeTransactionPurchase,
		uuid.NewString(),
		"123",
	)

	Ok(t, err)
	Equals(t, StatusPending, transaction.GetStatus())
}

func TestTransaction_WrongType(t *testing.T) {
	_, err := NewTransaction(
		uuid.NewString(),
		40,
		20,
		"",
		"send_email",
		uuid.NewString(),
		"123",
	)

	Equals(t, errors.INVALID_FIELD_VALUE("type").Error(), err.Error())
}

func TestTransaction_PurchaseWithoutAmount(t *testing.T) {
	_, err := NewTransaction(
		uuid.NewString(),
		40,
		0,
		"BRL",
		TypeTransactionPurchase,
		uuid.NewString(),
		"123",
	)

	Equals(t, errors.INVALID_FIELD_VALUE("type").Error(), err.Error())
}

func TestTransaction_PurchaseWithoutCurrency(t *testing.T) {
	_, err := NewTransaction(
		uuid.NewString(),
		40,
		20,
		"",
		TypeTransactionPurchase,
		uuid.NewString(),
		"123",
	)

	Equals(t, errors.INVALID_FIELD_VALUE("type").Error(), err.Error())
}

func TestTransaction_SetStatus(t *testing.T) {
	transaction, err := NewTransaction(
		uuid.NewString(),
		40,
		20,
		"BRL",
		TypeTransactionPurchase,
		uuid.NewString(),
		"123",
	)

	Ok(t, err)
	Equals(t, StatusPending, transaction.GetStatus())

	err = transaction.SetStatus(StatusCompleted)

	Ok(t, err)
	Equals(t, StatusCompleted, transaction.GetStatus())
}

func TestTransaction_SetStatusAfterCompleted(t *testing.T) {
	transaction, err := NewTransaction(
		uuid.NewString(),
		40,
		20,
		"BRL",
		TypeTransactionPurchase,
		uuid.NewString(),
		"123",
	)

	Ok(t, err)
	Equals(t, StatusPending, transaction.GetStatus())

	err = transaction.SetStatus(StatusCompleted)

	Ok(t, err)
	Equals(t, StatusCompleted, transaction.GetStatus())

	err = transaction.SetStatus(StatusPending)

	Equals(t, errors.FINISHED_OPERATION_ERROR().Error(), err.Error())
}

func TestTransaction_SetStatusAfterFailed(t *testing.T) {
	transaction, err := NewTransaction(
		uuid.NewString(),
		40,
		20,
		"BRL",
		TypeTransactionPurchase,
		uuid.NewString(),
		"123",
	)

	Ok(t, err)
	Equals(t, StatusPending, transaction.GetStatus())

	err = transaction.SetStatus(StatusFailed)

	Ok(t, err)
	Equals(t, StatusFailed, transaction.GetStatus())

	err = transaction.SetStatus(StatusCompleted)

	Equals(t, errors.FINISHED_OPERATION_ERROR().Error(), err.Error())
}
