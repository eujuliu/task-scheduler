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

var Client *Stripe

func Load(config *config.StripeConfig) *Stripe {
	if Client != nil {
		return Client
	}

	if config.APIKey == "" {
		panic(fmt.Errorf("stripe API key missing or not valid"))
	}

	sc := stripe.NewClient(config.APIKey)

	Client = &Stripe{
		client: sc,
		config: config,
	}

	return Client
}

func (s *Stripe) GetClient() *stripe.Client {
	return s.client
}
