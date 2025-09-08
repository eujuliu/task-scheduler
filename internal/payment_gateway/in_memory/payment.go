package in_memory_paymentgateway

import (
	"maps"

	"github.com/google/uuid"
)

type InMemoryPaymentPaymentGateway struct {
	payments []map[string]string
}

func NewInMemoryPaymentPaymentGateway() *InMemoryPaymentPaymentGateway {
	return &InMemoryPaymentPaymentGateway{
		payments: []map[string]string{},
	}
}

func (pg *InMemoryPaymentPaymentGateway) Create(
	internalID string,
	customerID *string,
	amount int,
	currency, idempotencyKey string,
	props *map[string]string,
) (string, error) {
	payment := map[string]string{
		"id":             uuid.NewString(),
		"internalID":     internalID,
		"customerID":     *customerID,
		"idempotencyKey": idempotencyKey,
	}

	if props != nil {
		maps.Copy(payment, *props)
	}

	pg.payments = append(pg.payments, payment)

	return payment["id"], nil
}
