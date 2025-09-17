package entities

import (
	"fmt"
	"scheduler/internal/errors"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	BaseEntity
	userId        string
	credits       int
	amount        int
	currency      string
	status        string
	kind          string
	referenceId   string
	idempotentKey string
}

const (
	StatusPending   string = "PENDING"
	StatusCompleted string = "COMPLETED"
	StatusFailed    string = "FAILED"
	StatusCanceled  string = "CANCELED"
	StatusFrozen    string = "FROZEN"
)

const (
	TypeTransactionPurchase string = "PURCHASE"
	TypeTransactionTaskSend string = "TASK_SEND"
)

var AvailableStatusPerType = map[string][]string{
	TypeTransactionPurchase: {StatusCompleted, StatusFailed, StatusCanceled},
	TypeTransactionTaskSend: {
		StatusFrozen,
		StatusCompleted,
		StatusFailed,
		StatusCanceled,
	},
}

func NewTransaction(
	userId string,
	credits int,
	amount int,
	currency string,
	kind string,
	referenceId string,
	idempotencyKey string,
) (*Transaction, error) {
	if uuid.Validate(userId) != nil {
		return nil, errors.INVALID_FIELD_VALUE("user id", nil)
	}

	if credits < 10 || credits > 100 {
		reason := "the credits can only be a number between 10 and 100"
		return nil, errors.INVALID_FIELD_VALUE("credits", &reason)
	}

	transactionTypes := []string{
		TypeTransactionPurchase,
		TypeTransactionTaskSend,
	}

	if !slices.Contains(transactionTypes, kind) {
		reason := fmt.Sprintf(
			"you need to set one of these (%s, %s)",
			TypeTransactionPurchase,
			TypeTransactionTaskSend,
		)
		return nil, errors.INVALID_FIELD_VALUE("type", &reason)
	}

	if kind == TypeTransactionPurchase && (amount <= 0 || len(currency) == 0) {
		reason := "for purchases you need to set the amount and the currency"
		return nil, errors.INVALID_FIELD_VALUE("type", &reason)
	}

	transaction := &Transaction{
		BaseEntity:    *NewBaseEntity(),
		userId:        userId,
		kind:          kind,
		amount:        amount,
		currency:      currency,
		credits:       credits,
		status:        StatusPending,
		referenceId:   referenceId,
		idempotentKey: idempotencyKey,
	}

	return transaction, nil
}

func HydrateTransaction(
	id, userId string,
	credits, amount int,
	currency, status, kind, referenceID, idempotentKey string,
	createdAt, updatedAt time.Time,
) *Transaction {
	return &Transaction{
		BaseEntity: BaseEntity{
			id:        id,
			createdAt: createdAt,
			updatedAt: updatedAt,
		},

		userId:        userId,
		credits:       credits,
		amount:        amount,
		currency:      currency,
		status:        status,
		kind:          kind,
		referenceId:   referenceID,
		idempotentKey: idempotentKey,
	}
}

func (t *Transaction) GetUserId() string {
	return t.userId
}

func (t *Transaction) GetCredits() int {
	return t.credits
}

func (t *Transaction) GetAmount() int {
	return t.amount
}

func (t *Transaction) GetCurrency() string {
	return t.currency
}

func (t *Transaction) SetStatus(status string) error {
	if t.readonly() {
		return errors.FINISHED_OPERATION_ERROR()
	}

	formatted := strings.TrimSpace(strings.ToUpper(status))

	_, ok := AvailableStatusPerType[t.kind]

	if !ok {
		return errors.INVALID_FIELD_VALUE("type", nil)
	}

	if !slices.Contains(AvailableStatusPerType[t.kind], formatted) {
		reason := "see the available status for this kind of transaction"
		return errors.INVALID_FIELD_VALUE("status", &reason)
	}

	t.status = formatted

	return nil
}

func (t *Transaction) GetStatus() string {
	return t.status
}

func (t *Transaction) GetType() string {
	return t.kind
}

func (t *Transaction) SetReferenceId(referenceId string) error {
	if t.referenceId != "" {
		reason := "the reference id is already added"
		return errors.INVALID_FIELD_VALUE("reference id", &reason)
	}

	t.referenceId = referenceId

	return nil
}

func (t *Transaction) GetReferenceId() string {
	return t.referenceId
}

func (t *Transaction) GetIdempotencyKey() string {
	return t.idempotentKey
}

func (t *Transaction) readonly() bool {
	return slices.Contains(
		[]string{StatusCompleted, StatusFailed, StatusCanceled},
		t.status,
	)
}
