package in_memory_paymentgateway

import (
	"maps"

	"github.com/google/uuid"
)

type InMemoryCustomerPaymentGateway struct {
	customers []map[string]string
}

func NewInMemoryCustomerPaymentGateway() *InMemoryCustomerPaymentGateway {
	return &InMemoryCustomerPaymentGateway{
		customers: []map[string]string{},
	}
}

func (pg *InMemoryCustomerPaymentGateway) GetFirstByEmail(email string) (*string, error) {
	for _, v := range pg.customers {
		if value, ok := v["email"]; ok && value == email {
			id := v["id"]

			return &id, nil
		}
	}

	return nil, nil
}

func (pg *InMemoryCustomerPaymentGateway) Create(
	internalID, name, email string,
	props *map[string]string,
) (string, error) {
	customer := map[string]string{
		"id":         uuid.NewString(),
		"internalID": internalID,
		"name":       name,
		"email":      email,
	}

	if props != nil {
		maps.Copy(customer, *props)
	}

	pg.customers = append(pg.customers, customer)

	return customer["id"], nil
}
