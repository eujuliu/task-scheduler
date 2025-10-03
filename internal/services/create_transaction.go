package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/interfaces"
)

type CreateTransactionService struct {
	userRepository         interfaces.IUserRepository
	transactionRepository  interfaces.ITransactionRepository
	paymentPaymentGateway  interfaces.IPaymentPaymentGateway
	customerPaymentGateway interfaces.ICustomerPaymentGateway
}

func NewCreateTransactionService(
	userRepository interfaces.IUserRepository,
	transactionRepository interfaces.ITransactionRepository,
	customerPaymentGateway interfaces.ICustomerPaymentGateway,
	paymentPaymentGateway interfaces.IPaymentPaymentGateway,
) *CreateTransactionService {
	return &CreateTransactionService{
		userRepository:         userRepository,
		transactionRepository:  transactionRepository,
		customerPaymentGateway: customerPaymentGateway,
		paymentPaymentGateway:  paymentPaymentGateway,
	}
}

func (s *CreateTransactionService) Execute(
	userId string,
	credits int,
	currency string,
	kind string,
	referenceId string,
	idempotencyKey string,
) (*entities.Transaction, error) {
	slog.Info("create transaction service started...")
	slog.Debug(fmt.Sprint("input ", userId,
		credits,
		currency,
		kind,
		idempotencyKey))

	user, _ := s.userRepository.GetFirstById(userId)
	if user == nil {
		slog.Error("user don't exists")
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	if kind == entities.TypeTransactionPurchase {
		referenceId = ""
	}

	transaction, _ := s.transactionRepository.GetFirstByIdempotencyKey(idempotencyKey)

	if transaction != nil {
		slog.Error("transaction with this idempotency key already exists")
		return transaction, nil
	}

	amount, err := s.calculateAmount(credits, currency)
	if err != nil {
		return nil, err
	}

	transaction, err = entities.NewTransaction(
		user.GetId(),
		credits,
		amount,
		currency,
		kind,
		referenceId,
		idempotencyKey,
	)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	err = s.transactionRepository.Create(transaction)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	if transaction.GetType() == entities.TypeTransactionPurchase {
		slog.Info("creating payment into payment gateway...")
		customer, _ := s.customerPaymentGateway.GetFirstByEmail(user.GetEmail())

		payment, err := s.paymentPaymentGateway.Create(
			transaction.GetId(),
			customer,
			amount,
			currency,
			idempotencyKey,
			nil,
		)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		err = transaction.SetReferenceId(payment)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		err = s.transactionRepository.Update(transaction)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		slog.Info("finish creation of payment into payment gateway...")
	}

	slog.Info("create transaction service finished...")
	slog.Debug(fmt.Sprintf("returned transaction %+v", transaction))

	return transaction, nil
}

func (s *CreateTransactionService) calculateAmount(credits int, currency string) (int, error) {
	/*
		Obviously that if you want to use this in production
		you will need to get the conversions from some api, but this project is only for
		learn proposes then i will continue with this static function
	*/

	creditsInDollars := float32(credits / 10)
	conversions := map[string]float32{
		"USD":  1,
		"BRL":  5.42,
		"EUR":  0.85,
		"TASK": 1,
	}

	conversion, ok := conversions[currency]

	if !ok {
		reason := "the valid currencies are (BRL, USD, EUR)"
		return 0, errors.INVALID_FIELD_VALUE("currency", &reason)
	}

	quantity := (creditsInDollars * conversion) / conversions["USD"]

	return int(quantity * 100), nil
}
