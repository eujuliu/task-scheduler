package services

import (
	"fmt"
	"scheduler/internal/entities"
	"scheduler/internal/errors"

	repos "scheduler/internal/repositories"
)

type UpdateTaskTransactionService struct {
	userRepository        repos.IUserRepository
	transactionRepository repos.ITransactionRepository
	errorRepository       repos.IErrorRepository
}

func NewUpdateTaskTransactionService(
	userRepository repos.IUserRepository,
	transactionRepository repos.ITransactionRepository,
	errorRepository repos.IErrorRepository,
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
	transaction, _ := s.transactionRepository.GetFirstById(transactionId)

	if transaction == nil {
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	user, _ := s.userRepository.GetFirstById(transaction.GetUserId())

	if user == nil {
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	err := transaction.SetStatus(entities.StatusCanceled)
	if err != nil {
		return nil, err
	}

	err = s.transactionRepository.Update(transaction)
	if err != nil {
		return nil, err
	}

	err = user.RemoveFrozenCredits(int(transaction.GetCredits()), false)
	if err != nil {
		return nil, err
	}

	err = s.userRepository.Update(user)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *UpdateTaskTransactionService) Fail(
	transactionId string,
	refund bool,
	reason string,
) (*entities.Transaction, error) {
	transaction, _ := s.transactionRepository.GetFirstById(transactionId)

	if transaction == nil {
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	user, _ := s.userRepository.GetFirstById(transaction.GetUserId())

	if user == nil {
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	err := transaction.SetStatus(entities.StatusFailed)
	if err != nil {
		return nil, err
	}

	err = s.transactionRepository.Update(transaction)
	if err != nil {
		return nil, err
	}

	err = user.RemoveFrozenCredits(int(transaction.GetCredits()), refund)
	if err != nil {
		return nil, err
	}

	err = s.userRepository.Update(user)
	if err != nil {
		return nil, err
	}

	error := entities.NewError(
		transaction.GetId(),
		entities.TypeErrorTask,
		reason,
		user.GetId(),
		map[string]string{"refund": fmt.Sprint(refund)},
	)

	err = s.errorRepository.Create(error)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *UpdateTaskTransactionService) Frozen(
	transactionId string,
) (*entities.Transaction, error) {
	transaction, _ := s.transactionRepository.GetFirstById(transactionId)

	if transaction == nil {
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	user, _ := s.userRepository.GetFirstById(transaction.GetUserId())

	if user == nil {
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	err := transaction.SetStatus(entities.StatusFrozen)
	if err != nil {
		return nil, err
	}

	err = s.transactionRepository.Update(transaction)
	if err != nil {
		return nil, err
	}

	err = user.AddFrozenCredits(int(transaction.GetCredits()))
	if err != nil {
		return nil, err
	}

	err = s.userRepository.Update(user)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *UpdateTaskTransactionService) Cancel(
	transactionId string,
) (*entities.Transaction, error) {
	transaction, _ := s.transactionRepository.GetFirstById(transactionId)

	if transaction == nil {
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	user, _ := s.userRepository.GetFirstById(transaction.GetUserId())

	if user == nil {
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	err := transaction.SetStatus(entities.StatusCanceled)
	if err != nil {
		return nil, err
	}

	err = s.transactionRepository.Update(transaction)
	if err != nil {
		return nil, err
	}

	err = user.RemoveFrozenCredits(int(transaction.GetCredits()), true)
	if err != nil {
		return nil, err
	}

	err = s.userRepository.Update(user)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
