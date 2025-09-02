package services

import (
	"fmt"
	"log/slog"
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
	slog.Info(fmt.Sprintf("getting user transaction with id %s", transactionId))
	user, _ := s.userRepository.GetFirstById(userId)

	if user == nil {
		slog.Error("user don't exists")
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	transaction, _ := s.transactionRepository.GetFirstById(transactionId)

	if transaction == nil {
		slog.Error("transaction don't exists")
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	slog.Info("transaction returned...")
	slog.Debug(fmt.Sprintf("Transaction %+v", transaction))

	return transaction, nil
}
