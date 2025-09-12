package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/interfaces"
)

type UpdatePurchaseTransactionService struct {
	userRepository        interfaces.IUserRepository
	transactionRepository interfaces.ITransactionRepository
	errorRepository       interfaces.IErrorRepository
}

func NewUpdatePurchaseTransactionService(
	userRepository interfaces.IUserRepository,
	transactionRepository interfaces.ITransactionRepository,
	errorRepository interfaces.IErrorRepository,
) *UpdatePurchaseTransactionService {
	return &UpdatePurchaseTransactionService{
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
		errorRepository:       errorRepository,
	}
}

func (s *UpdatePurchaseTransactionService) Complete(
	transactionId string,
) (*entities.Transaction, error) {
	slog.Info("update purchase transaction (complete) service started...")
	slog.Debug(fmt.Sprint("input ", transactionId))

	transaction, _ := s.transactionRepository.GetFirstById(transactionId)

	if transaction == nil {
		slog.Error("transaction not found error")
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	user, _ := s.userRepository.GetFirstById(transaction.GetUserId())

	if user == nil {
		slog.Error("user not found error")
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	user.AddCredits(int(transaction.GetCredits()))

	err := s.userRepository.Update(user)
	if err != nil {
		slog.Error(fmt.Sprintf("update user error %s", err.Error()))
		return nil, err
	}

	err = s.updateTransactionStatus(transaction, entities.StatusCompleted)
	if err != nil {
		slog.Error(fmt.Sprintf("transaction status update error %s", err.Error()))
		return nil, err
	}

	slog.Info("update purchase transaction (complete) service finished...")
	slog.Debug(fmt.Sprintf("returned transaction %+v", transaction))

	return transaction, nil
}

func (s *UpdatePurchaseTransactionService) Fail(
	transactionId string,
	reason string,
) (*entities.Transaction, error) {
	slog.Info("update purchase transaction (fail) service started...")
	slog.Debug(fmt.Sprint("input ", transactionId))

	transaction, _ := s.transactionRepository.GetFirstById(transactionId)

	if transaction == nil {
		slog.Error("transaction not found error")
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	user, _ := s.userRepository.GetFirstById(transaction.GetUserId())

	if user == nil {
		slog.Error("user not found error")
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	err := s.updateTransactionStatus(transaction, entities.StatusFailed)
	if err != nil {
		slog.Error(fmt.Sprintf("transaction status update error %s", err.Error()))
		return nil, err
	}

	error := entities.NewError(
		transaction.GetId(),
		entities.TypeErrorTransaction,
		reason,
		user.GetId(),
	)

	err = s.errorRepository.Create(error)
	if err != nil {
		return nil, err
	}

	slog.Info("update purchase transaction (fail) service finished...")
	slog.Debug(fmt.Sprintf("returned transaction %+v", transaction))
	return transaction, nil
}

func (s *UpdatePurchaseTransactionService) updateTransactionStatus(
	transaction *entities.Transaction,
	status string,
) error {
	err := transaction.SetStatus(status)
	if err != nil {
		return err
	}

	err = s.transactionRepository.Update(transaction)
	if err != nil {
		return err
	}

	return nil
}
