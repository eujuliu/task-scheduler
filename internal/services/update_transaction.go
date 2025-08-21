package services

import (
	"fmt"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	repos "scheduler/internal/repositories"
)

type UpdateTransactionService struct {
	userRepository        repos.IUserRepository
	transactionRepository repos.ITransactionRepository
	errorRepository       repos.IErrorRepository
}

func NewUpdateTransactionService(userRepository repos.IUserRepository, transactionRepository repos.ITransactionRepository, errorRepository repos.IErrorRepository) *UpdateTransactionService {
	return &UpdateTransactionService{
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
		errorRepository:       errorRepository,
	}
}

func (s *UpdateTransactionService) Execute(transactionId string, newStatus string, options map[string]any) (*entities.Transaction, error) {
	transaction, err := s.transactionRepository.GetFirstById(transactionId)

	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetFirstById(transaction.GetUserId())

	if err != nil {
		return nil, err
	}

	if _, ok := options["reason"]; !ok && newStatus == "failed" {
		return nil, errors.MISSING_PARAM_ERROR("option.reason")
	}

	if _, ok := options["refund"]; !ok && newStatus == "failed" && transaction.GetType() != "purchase" {
		return nil, errors.MISSING_PARAM_ERROR("option.refund")
	}

	switch transaction.GetType() {
	case "purchase":
		return s.updatePurchase(newStatus, transaction, user, options)
	case "task_send":
		return s.updateTaskSend(newStatus, transaction, user, options)
	default:
		return nil, errors.INVALID_FIELD_VALUE("transaction type")
	}
}

func (s *UpdateTransactionService) updatePurchase(newStatus string, transaction *entities.Transaction, user *entities.User, options map[string]any) (*entities.Transaction, error) {
	if newStatus == "completed" {
		err := transaction.SetStatus(newStatus)

		if err != nil {
			return nil, err
		}

		err = s.transactionRepository.Update(transaction)

		if err != nil {
			return nil, err
		}

		user.AddCredits(transaction.GetCredits())

		err = s.userRepository.Update(user)

		if err != nil {
			return nil, err
		}

		return transaction, nil
	}

	if newStatus == "failed" {
		err := transaction.SetStatus(newStatus)

		if err != nil {
			return nil, err
		}

		err = s.transactionRepository.Update(transaction)

		if err != nil {
			return nil, err
		}

		formattedMap := make(map[string]string)

		for k, v := range options {
			formattedMap[k] = fmt.Sprint(v)
		}

		error := entities.NewError(transaction.GetId(), "transaction_purchase", formattedMap["reason"], user.GetId(), formattedMap)

		err = s.errorRepository.Create(error)

		if err != nil {
			return nil, err
		}

		return transaction, nil
	}

	return nil, errors.INVALID_FIELD_VALUE("status")
}

func (s *UpdateTransactionService) updateTaskSend(newStatus string, transaction *entities.Transaction, user *entities.User, options map[string]any) (*entities.Transaction, error) {
	if newStatus == "frozen" {
		err := transaction.SetStatus(newStatus)

		if err != nil {
			return nil, err
		}

		err = s.transactionRepository.Update(transaction)

		if err != nil {
			return nil, err
		}

		err = user.AddFrozenCredits(transaction.GetCredits())

		if err != nil {
			return nil, err
		}

		err = s.userRepository.Update(user)

		if err != nil {
			return nil, err
		}

		return transaction, nil
	}

	if newStatus == "completed" {
		err := transaction.SetStatus(newStatus)

		if err != nil {
			return nil, err
		}

		err = s.transactionRepository.Update(transaction)

		if err != nil {
			return nil, err
		}

		err = user.RemoveFrozenCredits(transaction.GetCredits(), false)

		if err != nil {
			return nil, err
		}

		err = s.userRepository.Update(user)

		if err != nil {
			return nil, err
		}

		return transaction, nil
	}

	if newStatus == "failed" {
		refund, ok := options["refund"].(bool)

		if !ok {
			return nil, errors.INVALID_TYPE()
		}

		err := transaction.SetStatus(newStatus)

		if err != nil {
			return nil, err
		}

		err = s.transactionRepository.Update(transaction)

		if err != nil {
			return nil, err
		}

		err = user.RemoveFrozenCredits(transaction.GetCredits(), refund)

		if err != nil {
			return nil, err
		}

		err = s.userRepository.Update(user)

		if err != nil {
			return nil, err
		}

		formattedMap := make(map[string]string)

		for k, v := range options {
			formattedMap[k] = fmt.Sprint(v)
		}

		error := entities.NewError(transaction.GetId(), "transaction_task_send", formattedMap["reason"], user.GetId(), formattedMap)

		err = s.errorRepository.Create(error)

		if err != nil {
			return nil, err
		}

		return transaction, nil
	}

	return nil, errors.INVALID_FIELD_VALUE("status")
}
