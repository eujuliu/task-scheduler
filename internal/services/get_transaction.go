package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/interfaces"
)

type GetTransactionService struct {
	userRepository        interfaces.IUserRepository
	transactionRepository interfaces.ITransactionRepository
}

func NewGetTransactionService(
	userRepository interfaces.IUserRepository,
	transactionRepository interfaces.ITransactionRepository,
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
	slog.Info("get transaction service started...")
	slog.Debug(fmt.Sprint("input ", userId, transactionId))
	user, _ := s.userRepository.GetFirstById(userId)

	if user == nil {
		slog.Error("user not exists error")
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	transaction, _ := s.transactionRepository.GetFirstById(transactionId)

	if transaction == nil {
		slog.Error("transaction not exists error")
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	slog.Info("get transaction service finished...")
	slog.Debug(fmt.Sprintf("transaction returned %+v", transaction))

	return transaction, nil
}
