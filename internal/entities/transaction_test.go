package entities_test

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	. "scheduler/test"
	"testing"

	"github.com/google/uuid"
)

func TestTransaction_New(t *testing.T) {
	transaction, err := entities.NewTransaction(uuid.NewString(), 40, 20, "BRL", "purchase", uuid.NewString(), "123")

	Ok(t, err)
	Equals(t, "frozen", transaction.GetStatus())
}

func TestTransaction_WrongType(t *testing.T) {
	_, err := entities.NewTransaction(uuid.NewString(), 40, 20, "", "send_email", uuid.NewString(), "123")

	Equals(t, errors.INVALID_FIELD_VALUE("type").Error(), err.Error())
}

func TestTransaction_PurchaseWithoutAmount(t *testing.T) {
	_, err := entities.NewTransaction(uuid.NewString(), 40, 0, "BRL", "purchase", uuid.NewString(), "123")

	Equals(t, errors.INVALID_FIELD_VALUE("type").Error(), err.Error())
}

func TestTransaction_PurchaseWithoutCurrency(t *testing.T) {
	_, err := entities.NewTransaction(uuid.NewString(), 40, 20, "", "purchase", uuid.NewString(), "123")

	Equals(t, errors.INVALID_FIELD_VALUE("type").Error(), err.Error())
}

func TestTransaction_SetStatus(t *testing.T) {
	transaction, err := entities.NewTransaction(uuid.NewString(), 40, 20, "BRL", "purchase", uuid.NewString(), "123")

	Ok(t, err)
	Equals(t, "frozen", transaction.GetStatus())

	err = transaction.SetStatus("completed")

	Ok(t, err)
	Equals(t, "completed", transaction.GetStatus())
}

func TestTransaction_SetStatusAfterCompleted(t *testing.T) {
	transaction, err := entities.NewTransaction(uuid.NewString(), 40, 20, "BRL", "purchase", uuid.NewString(), "123")

	Ok(t, err)
	Equals(t, "frozen", transaction.GetStatus())

	err = transaction.SetStatus("completed")

	Ok(t, err)
	Equals(t, "completed", transaction.GetStatus())

	err = transaction.SetStatus("frozen")

	Equals(t, errors.FINISHED_OPERATION_ERROR().Error(), err.Error())
}

func TestTransaction_SetStatusAfterFailed(t *testing.T) {
	transaction, err := entities.NewTransaction(uuid.NewString(), 40, 20, "BRL", "purchase", uuid.NewString(), "123")

	Ok(t, err)
	Equals(t, "frozen", transaction.GetStatus())

	err = transaction.SetStatus("failed")

	Ok(t, err)
	Equals(t, "failed", transaction.GetStatus())

	err = transaction.SetStatus("completed")

	Equals(t, errors.FINISHED_OPERATION_ERROR().Error(), err.Error())
}
