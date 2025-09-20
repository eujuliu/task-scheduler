package stripe_paymentgateway

import (
	"context"
	"maps"
	"scheduler/internal/errors"
	"scheduler/pkg/stripe"

	stripeGo "github.com/stripe/stripe-go/v82"
)

type StripePaymentPaymentGateway struct {
	stripe *stripe.Stripe
}

func NewStripePaymentPaymentGateway(stripe *stripe.Stripe) *StripePaymentPaymentGateway {
	return &StripePaymentPaymentGateway{
		stripe: stripe,
	}
}

func (pg *StripePaymentPaymentGateway) Create(
	internalID string,
	customerID *string,
	amount int,
	currency, idempotencyKey string,
	props *map[string]string,
) (string, error) {
	allowedCurrencies := map[string]stripeGo.Currency{
		"USD": stripeGo.CurrencyUSD, "BRL": stripeGo.CurrencyBRL, "EUR": stripeGo.CurrencyEUR,
	}

	metadata := map[string]string{
		"transactionID": internalID,
	}

	if props != nil {
		maps.Copy(*props, metadata)
	}

	if _, ok := allowedCurrencies[currency]; !ok {
		reason := "you can only use these (USD, BRL, EUR)"
		return "", errors.INVALID_FIELD_VALUE("currency", &reason)
	}

	params := &stripeGo.PaymentIntentCreateParams{
		Amount:   stripeGo.Int64(int64(amount)),
		Currency: stripeGo.String(allowedCurrencies[currency]),
		Customer: customerID,
		Metadata: metadata,
	}

	params.SetIdempotencyKey(idempotencyKey)

	result, err := pg.stripe.Client().V1PaymentIntents.Create(context.TODO(), params)

	return result.ClientSecret, err
}
