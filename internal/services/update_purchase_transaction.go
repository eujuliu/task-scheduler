package services

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	repos "scheduler/internal/repositories"
)

type UpdatePurchaseTransactionService struct {
	userRepository        repos.IUserRepository
	transactionRepository repos.ITransactionRepository
	errorRepository       repos.IErrorRepository
}

func NewUpdatePurchaseTransactionService(
	userRepository repos.IUserRepository,
	transactionRepository repos.ITransactionRepository,
	errorRepository repos.IErrorRepository,
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
	transaction, _ := s.transactionRepository.GetFirstById(transactionId)

	if transaction == nil {
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	user, _ := s.userRepository.GetFirstById(transaction.GetUserId())

	if user == nil {
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	user.AddCredits(transaction.GetCredits())

	err := transaction.SetStatus(entities.StatusCompleted)
	if err != nil {
		return nil, err
	}

	err = s.transactionRepository.Update(transaction)
	if err != nil {
		return nil, err
	}

	err = s.userRepository.Update(user)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *UpdatePurchaseTransactionService) Fail(
	transactionId string,
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

	error := entities.NewError(
		transaction.GetId(),
		entities.TypeErrorTransaction,
		reason,
		user.GetId(),
		make(map[string]string),
	)

	err = s.errorRepository.Create(error)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
