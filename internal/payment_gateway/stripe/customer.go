package stripe_paymentgateway

import (
	"context"
	"fmt"
	"maps"
	"scheduler/pkg/stripe"

	stripeGo "github.com/stripe/stripe-go/v82"
)

type StripeCustomerPaymentGateway struct {
	instance *stripe.Stripe
}

func NewStripeCustomerPaymentGateway() *StripeCustomerPaymentGateway {
	return &StripeCustomerPaymentGateway{
		instance: stripe.Client,
	}
}

func (pg *StripeCustomerPaymentGateway) GetFirstByEmail(email string) (*string, error) {
	params := &stripeGo.CustomerSearchParams{
		SearchParams: stripeGo.SearchParams{
			Query: fmt.Sprintf("email:'%s'", email),
		},
	}
	result := pg.instance.GetClient().V1Customers.Search(context.TODO(), params)

	for customer, err := range result {
		if err != nil {
			return nil, err
		}

		return &customer.ID, nil
	}

	return nil, nil
}

func (pg *StripeCustomerPaymentGateway) Create(
	internalID, name, email string,
	props *map[string]string,
) (string, error) {
	metadata := map[string]string{
		"UserID": internalID,
	}

	if props != nil {
		maps.Copy(*props, metadata)
	}

	params := &stripeGo.CustomerCreateParams{
		Email:    stripeGo.String(email),
		Name:     stripeGo.String(name),
		Metadata: metadata,
	}

	result, err := pg.instance.GetClient().V1Customers.Create(context.TODO(), params)

	return result.ID, err
}
