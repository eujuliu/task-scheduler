package services

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	repos "scheduler/internal/repositories"
)

type CreateTransactionService struct {
	userRepository        repos.IUserRepository
	transactionRepository repos.ITransactionRepository
}

func NewCreateTransactionService(userRepository repos.IUserRepository, transactionRepository repos.ITransactionRepository) *CreateTransactionService {
	return &CreateTransactionService{
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
	}
}

func (s *CreateTransactionService) Execute(userId string, credits int, amount int, currency string, kind string, referenceId string, idempotencyKey string) (*entities.Transaction, error) {
	user, err := s.userRepository.GetFirstById(userId)

	if err != nil {
		return nil, err
	}

	exists, _ := s.transactionRepository.GetFirstByReferenceId(referenceId)

	if exists != nil {
		return nil, errors.TRANSACTION_ALREADY_EXISTS_ERROR()
	}

	exists, _ = s.transactionRepository.GetFirstByIdempotencyKey(idempotencyKey)

	if exists != nil {
		return nil, errors.TRANSACTION_ALREADY_EXISTS_ERROR()
	}

	transaction, err := entities.NewTransaction(user.GetId(), credits, amount, currency, kind, referenceId, idempotencyKey)

	if err != nil {
		return nil, err
	}

	err = s.transactionRepository.Create(transaction)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}
