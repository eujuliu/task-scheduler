package paymentgateway

type ICustomerPaymentGateway interface {
	GetFirstByEmail(email string) (*string, error)
	Create(internalID, username, email string, props *map[string]string) (string, error)
}

type IPaymentPaymentGateway interface {
	Create(
		internalID string,
		customerID *string,
		amount int,
		currency, idempotencyKey string,
		props *map[string]string,
	) (string, error)
}
