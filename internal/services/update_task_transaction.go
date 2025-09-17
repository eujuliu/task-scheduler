package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/interfaces"
)

type UpdateTaskTransactionService struct {
	userRepository        interfaces.IUserRepository
	transactionRepository interfaces.ITransactionRepository
	errorRepository       interfaces.IErrorRepository
}

func NewUpdateTaskTransactionService(
	userRepository interfaces.IUserRepository,
	transactionRepository interfaces.ITransactionRepository,
	errorRepository interfaces.IErrorRepository,
) *UpdateTaskTransactionService {
	return &UpdateTaskTransactionService{
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
		errorRepository:       errorRepository,
	}
}

func (s *UpdateTaskTransactionService) Complete(
	transactionId string,
) (*entities.Transaction, error) {
	slog.Info("update task transaction (complete) service started...")
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

	err := s.updateTransactionStatus(transaction, entities.StatusCompleted)
	if err != nil {
		slog.Error(fmt.Sprintf("transaction status update error %s", err.Error()))
		return nil, err
	}

	err = s.removeUserFrozenCredits(transaction, user, false)
	if err != nil {
		slog.Error(fmt.Sprintf("removing user frozen credits error %s", err.Error()))
		return nil, err
	}

	slog.Info("update task transaction (complete) service finished...")
	slog.Debug(fmt.Sprintf("returned transaction %+v", transaction))
	return transaction, nil
}

func (s *UpdateTaskTransactionService) Fail(
	transactionId string,
	refund bool,
	reason string,
) (*entities.Transaction, error) {
	slog.Info("update task transaction (fail) service started...")
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

	err = s.removeUserFrozenCredits(transaction, user, refund)
	if err != nil {
		slog.Error(fmt.Sprintf("removing user frozen credits error %s", err.Error()))
		return nil, err
	}

	error := entities.NewError(
		transaction.GetId(),
		entities.TypeErrorTask,
		reason,
		user.GetId(),
	)

	err = s.errorRepository.Create(error)
	if err != nil {

		slog.Error(fmt.Sprintf("create error entity error %s", err.Error()))
		return nil, err
	}

	slog.Info("update task transaction (fail) service finished...")
	slog.Debug(fmt.Sprintf("returned transaction %+v", transaction))
	return transaction, nil
}

func (s *UpdateTaskTransactionService) Frozen(
	transactionId string,
) (*entities.Transaction, error) {
	slog.Info("update task transaction (frozen) service started...")
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

	err := transaction.SetStatus(entities.StatusFrozen)
	if err != nil {
		slog.Error(fmt.Sprintf("transaction set status error %s", err.Error()))
		return nil, err
	}

	err = s.transactionRepository.Update(transaction)
	if err != nil {
		slog.Error(fmt.Sprintf("update transaction error error %s", err.Error()))
		return nil, err
	}

	err = user.AddFrozenCredits(transaction.GetCredits())
	if err != nil {
		slog.Error(fmt.Sprintf("remove frozen credits error %s", err.Error()))
		return nil, err
	}

	err = s.userRepository.Update(user)
	if err != nil {
		slog.Error(fmt.Sprintf("user update error %s", err.Error()))
		return nil, err
	}

	slog.Info("update task transaction (frozen) service finished...")
	slog.Debug(fmt.Sprintf("returned transaction %+v", transaction))
	return transaction, nil
}

func (s *UpdateTaskTransactionService) Cancel(
	transactionId string,
) (*entities.Transaction, error) {
	slog.Info("update task transaction (cancel) service started...")
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

	err := s.updateTransactionStatus(transaction, entities.StatusCanceled)
	if err != nil {
		slog.Error(fmt.Sprintf("transaction status update error %s", err.Error()))
		return nil, err
	}

	err = s.removeUserFrozenCredits(transaction, user, true)
	if err != nil {
		slog.Error(fmt.Sprintf("removing user frozen credits error %s", err.Error()))
		return nil, err
	}

	slog.Info("update task transaction (cancel) service finished...")
	slog.Debug(fmt.Sprintf("returned transaction %+v", transaction))

	return transaction, nil
}

func (s *UpdateTaskTransactionService) updateTransactionStatus(
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

func (s *UpdateTaskTransactionService) removeUserFrozenCredits(
	transaction *entities.Transaction,
	user *entities.User,
	refund bool,
) error {
	err := user.RemoveFrozenCredits(transaction.GetCredits(), refund)
	if err != nil {
		slog.Error(fmt.Sprintf("remove frozen credits error %s", err.Error()))
		return err
	}

	err = s.userRepository.Update(user)

	slog.Debug(fmt.Sprintf("user %+v", user))
	if err != nil {
		slog.Error(fmt.Sprintf("user update error %s", err.Error()))
		return err
	}

	return nil
}
