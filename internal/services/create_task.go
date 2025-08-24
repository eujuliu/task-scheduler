package services

import (
	"fmt"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	repos "scheduler/internal/repositories"
	"time"
)

type CreateTaskService struct {
	userRepository        repos.IUserRepository
	transactionRepository repos.ITransactionRepository
	taskRepository        repos.ITaskRepository
}

func NewCreateTaskService(userRepository repos.IUserRepository, transactionRepository repos.ITransactionRepository, taskRepository repos.ITaskRepository) *CreateTaskService {
	return &CreateTaskService{
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
		taskRepository:        taskRepository,
	}
}

func (s *CreateTaskService) Execute(kind string, runAt time.Time, timezone string, priority int, userId string, referenceId string, idempotencyKey string) (*entities.Task, error) {
	task, _ := s.taskRepository.GetFirstByReferenceId(referenceId)

	if task != nil {
		return nil, errors.TASK_ALREADY_EXISTS_ERROR()
	}

	task, _ = s.taskRepository.GetFirstByIdempotencyKey(idempotencyKey)

	if task != nil {
		return nil, errors.TASK_ALREADY_EXISTS_ERROR()
	}

	user, err := s.userRepository.GetFirstById(userId)

	if err != nil {
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	task, err = entities.NewTask(kind, user.GetId(), 10, runAt, timezone, referenceId, idempotencyKey)

	if err != nil {
		return nil, err
	}

	if priority != entities.PriorityLow {
		err = task.SetPriority(priority)

		if err != nil {
			return nil, err
		}
	}

	err = s.taskRepository.Create(task)

	if err != nil {
		return nil, err
	}

	transaction, _ := s.transactionRepository.GetFirstByReferenceId(referenceId)

	if transaction != nil {
		return nil, errors.TRANSACTION_ALREADY_EXISTS_ERROR()
	}

	transaction, _ = s.transactionRepository.GetFirstByIdempotencyKey(idempotencyKey)

	if transaction != nil {
		return nil, errors.TRANSACTION_ALREADY_EXISTS_ERROR()
	}

	transaction, err = entities.NewTransaction(user.GetId(), 10, 0, "", entities.TypeTransactionTaskSend, fmt.Sprintf("task_%s", referenceId), idempotencyKey)

	if err != nil {
		return nil, err
	}

	err = transaction.SetStatus(entities.StatusFrozen)

	if err != nil {
		return nil, err
	}

	err = s.transactionRepository.Create(transaction)

	if err != nil {
		return nil, err
	}

	return task, nil
}
