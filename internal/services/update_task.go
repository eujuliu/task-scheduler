package services

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	repos "scheduler/internal/repositories"
	"time"
)

type UpdateTaskService struct {
	taskRepository               repos.ITaskRepository
	transactionRepository        repos.ITransactionRepository
	updateTaskTransactionService *UpdateTaskTransactionService
}

func NewUpdateTaskService(
	taskRepository repos.ITaskRepository,
	transactionRepository repos.ITransactionRepository,
	updateTaskTransaction *UpdateTaskTransactionService,
) *UpdateTaskService {
	return &UpdateTaskService{
		taskRepository:               taskRepository,
		transactionRepository:        transactionRepository,
		updateTaskTransactionService: updateTaskTransaction,
	}
}

func (s *UpdateTaskService) Execute(
	taskId string,
	runAt time.Time,
	timezone string,
	priority int,
) (*entities.Task, error) {
	task, _ := s.taskRepository.GetFirstById(taskId)

	if task == nil {
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	err := task.SetRunAt(runAt)
	if err != nil {
		return nil, err
	}

	err = task.SetTimezone(timezone)
	if err != nil {
		return nil, err
	}

	err = task.SetPriority(priority)
	if err != nil {
		return nil, err
	}

	err = s.taskRepository.Update(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *UpdateTaskService) Complete(
	taskId string,
) (*entities.Task, error) {
	task, _ := s.taskRepository.GetFirstById(taskId)

	if task == nil {
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	err := task.SetStatus(entities.StatusCompleted)
	if err != nil {
		return nil, err
	}

	transaction, _ := s.transactionRepository.GetFirstByReferenceId(task.GetReferenceId())

	if transaction == nil {
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	_, err = s.updateTaskTransactionService.Complete(transaction.GetId())
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *UpdateTaskService) Fail(
	taskId string,
	refund bool,
	reason string,
) (*entities.Task, error) {
	task, _ := s.taskRepository.GetFirstById(taskId)

	if task == nil {
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	err := task.SetStatus(entities.StatusFailed)
	if err != nil {
		return nil, err
	}

	transaction, _ := s.transactionRepository.GetFirstByReferenceId(task.GetReferenceId())

	if transaction == nil {
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	_, err = s.updateTaskTransactionService.Fail(transaction.GetId(), refund, reason)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *UpdateTaskService) Cancel(
	taskId string,
) (*entities.Task, error) {
	task, _ := s.taskRepository.GetFirstById(taskId)

	if task == nil {
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	err := task.SetStatus(entities.StatusCanceled)
	if err != nil {
		return nil, err
	}

	transaction, _ := s.transactionRepository.GetFirstByReferenceId(task.GetReferenceId())

	if transaction == nil {
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	_, err = s.updateTaskTransactionService.Cancel(transaction.GetId())
	if err != nil {
		return nil, err
	}

	return task, nil
}
