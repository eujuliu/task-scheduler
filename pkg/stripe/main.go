package stripe

import (
	"fmt"
	"scheduler/internal/config"

	"github.com/stripe/stripe-go/v82"
)

type Stripe struct {
	client *stripe.Client
	config *config.StripeConfig
}

func NewStripe(config *config.StripeConfig) (*Stripe, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("stripe API key missing or not valid")
	}

	sc := stripe.NewClient(config.APIKey)

	return &Stripe{
		client: sc,
		config: config,
	}, nil
}

func (s *Stripe) Client() *stripe.Client {
	return s.client
}
