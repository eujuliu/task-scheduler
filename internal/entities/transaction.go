package entities

import (
	"scheduler/internal/errors"
	"slices"
	"strings"

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
	StatusRunning   string = "RUNNING"
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
	TypeTransactionTaskSend: {StatusFrozen, StatusCompleted, StatusFailed, StatusCanceled},
}

func NewTransaction(userId string, credits int, amount int, currency string, kind string, referenceId string, idempotencyKey string) (*Transaction, error) {
	if uuid.Validate(userId) != nil {
		return nil, errors.INVALID_FIELD_VALUE("user id")
	}

	var transactionTypes = []string{TypeTransactionPurchase, TypeTransactionTaskSend}

	if !slices.Contains(transactionTypes, kind) {
		return nil, errors.INVALID_FIELD_VALUE("type")
	}

	if kind == TypeTransactionPurchase && (amount <= 0 || len(currency) == 0) {
		return nil, errors.INVALID_FIELD_VALUE("type")
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
		return errors.INVALID_FIELD_VALUE("type")
	}

	if !slices.Contains(AvailableStatusPerType[t.kind], formatted) {
		return errors.INVALID_FIELD_VALUE("status")
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

func (t *Transaction) GetReferenceId() string {
	return t.referenceId
}

func (t *Transaction) GetIdempotencyKey() string {
	return t.idempotentKey
}

func (t *Transaction) readonly() bool {
	return slices.Contains([]string{StatusCompleted, StatusFailed, StatusCanceled}, t.status)
}
