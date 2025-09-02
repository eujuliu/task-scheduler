package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	repos "scheduler/internal/repositories"
)

type CreateTransactionService struct {
	userRepository        repos.IUserRepository
	transactionRepository repos.ITransactionRepository
}

func NewCreateTransactionService(
	userRepository repos.IUserRepository,
	transactionRepository repos.ITransactionRepository,
) *CreateTransactionService {
	return &CreateTransactionService{
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
	}
}

func (s *CreateTransactionService) Execute(
	userId string,
	credits int,
	amount int,
	currency string,
	kind string,
	referenceId string,
	idempotencyKey string,
) (*entities.Transaction, error) {
	slog.Info("create transaction service started...")

	user, _ := s.userRepository.GetFirstById(userId)
	if user == nil {
		slog.Error("user don't exists")
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	exists, _ := s.transactionRepository.GetFirstByReferenceId(referenceId)

	if exists != nil {
		slog.Error("transaction with this reference id already exists")
		return nil, errors.TRANSACTION_ALREADY_EXISTS_ERROR()
	}

	exists, _ = s.transactionRepository.GetFirstByIdempotencyKey(idempotencyKey)

	if exists != nil {
		slog.Error("transaction with this idempotency key already exists")
		return nil, errors.TRANSACTION_ALREADY_EXISTS_ERROR()
	}

	transaction, err := entities.NewTransaction(
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

	slog.Info("create transaction service finished...")
	slog.Debug(fmt.Sprintf("returned transaction %+v", transaction))

	return transaction, nil
}
