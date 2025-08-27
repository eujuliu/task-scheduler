package services

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	repos "scheduler/internal/repositories"
)

type GetTransactionService struct {
	userRepository        repos.IUserRepository
	transactionRepository repos.ITransactionRepository
}

func NewGetTransactionService(
	userRepository repos.IUserRepository,
	transactionRepository repos.ITransactionRepository,
) *GetTransactionService {
	return &GetTransactionService{
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
	}
}

func (s *GetTransactionService) Execute(
	userId string,
	transactionId string,
) (*entities.Transaction, error) {
	user, _ := s.userRepository.GetFirstById(userId)

	if user == nil {
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	transaction, _ := s.transactionRepository.GetFirstById(transactionId)

	if transaction == nil {
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	return transaction, nil
}
