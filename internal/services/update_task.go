package services

import (
	"fmt"
	"log/slog"
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
	slog.Info("update task service started...")
	slog.Debug(fmt.Sprint("input ", taskId, runAt, timezone, priority))
	task, _ := s.taskRepository.GetFirstById(taskId)

	if task == nil {
		slog.Error("task not found error")
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	err := task.SetRunAt(runAt)
	if err != nil {
		slog.Error(fmt.Sprintf("set run at error %s", err.Error()))
		return nil, err
	}

	err = task.SetTimezone(timezone)
	if err != nil {
		slog.Error(fmt.Sprintf("set timezone error %s", err.Error()))
		return nil, err
	}

	err = task.SetPriority(priority)
	if err != nil {
		slog.Error(fmt.Sprintf("set priority error %s", err.Error()))
		return nil, err
	}

	err = s.taskRepository.Update(task)
	if err != nil {
		slog.Error(fmt.Sprintf("task update error %s", err.Error()))
		return nil, err
	}

	slog.Info("update task service finished...")
	slog.Debug(fmt.Sprintf("returned task %+v", task))

	return task, nil
}

func (s *UpdateTaskService) Complete(
	taskId string,
) (*entities.Task, error) {
	slog.Info("update task (complete) service started...")
	slog.Debug(fmt.Sprint("input ", taskId))
	task, _ := s.taskRepository.GetFirstById(taskId)

	if task == nil {
		slog.Error("task not found error")
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	err := s.updateTaskStatus(task, entities.StatusCompleted)
	if err != nil {
		slog.Error(fmt.Sprintf("task status update error %s", err.Error()))
		return nil, err
	}

	transaction, _ := s.transactionRepository.GetFirstByReferenceId(task.GetReferenceId())

	if transaction == nil {
		slog.Error("transaction not found error")
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	_, err = s.updateTaskTransactionService.Complete(transaction.GetId())
	if err != nil {

		slog.Error(fmt.Sprintf("update transaction error %s", err.Error()))
		return nil, err
	}

	slog.Info("update task (complete) service finished...")
	slog.Debug(fmt.Sprintf("returned task %+v", task))
	return task, nil
}

func (s *UpdateTaskService) Fail(
	taskId string,
	refund bool,
	reason string,
) (*entities.Task, error) {
	slog.Info("update task (fail) service started...")
	slog.Debug(fmt.Sprint("input ", taskId))

	task, _ := s.taskRepository.GetFirstById(taskId)

	if task == nil {
		slog.Error("task not found error")
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	err := s.updateTaskStatus(task, entities.StatusFailed)
	if err != nil {
		slog.Error(fmt.Sprintf("task status update error %s", err.Error()))
		return nil, err
	}

	transaction, _ := s.transactionRepository.GetFirstByReferenceId(task.GetReferenceId())

	if transaction == nil {
		slog.Error("transaction not found error")
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	_, err = s.updateTaskTransactionService.Fail(transaction.GetId(), refund, reason)
	if err != nil {
		slog.Error(fmt.Sprintf("update transaction error %s", err.Error()))
		return nil, err
	}

	slog.Info("update task service (fail) finished...")
	slog.Debug(fmt.Sprintf("returned task %+v", task))
	return task, nil
}

func (s *UpdateTaskService) Cancel(
	taskId string,
) (*entities.Task, error) {
	slog.Info("update task (cancel) service started...")
	slog.Debug(fmt.Sprint("input ", taskId))

	task, _ := s.taskRepository.GetFirstById(taskId)

	if task == nil {
		slog.Error("task not found error")
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	err := s.updateTaskStatus(task, entities.StatusCanceled)
	if err != nil {
		slog.Error(fmt.Sprintf("task status update error %s", err.Error()))
		return nil, err
	}

	transaction, _ := s.transactionRepository.GetFirstByReferenceId(task.GetReferenceId())

	if transaction == nil {
		slog.Error("transaction not found error")
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	_, err = s.updateTaskTransactionService.Cancel(transaction.GetId())
	if err != nil {
		slog.Error(fmt.Sprintf("update transaction error %s", err.Error()))
		return nil, err
	}

	slog.Info("update task service (cancel) finished...")
	slog.Debug(fmt.Sprintf("returned task %+v", task))
	return task, nil
}

func (s *UpdateTaskService) updateTaskStatus(task *entities.Task, status string) error {
	err := task.SetStatus(status)
	if err != nil {
		return err
	}

	err = s.taskRepository.Update(task)
	if err != nil {
		return err
	}

	return nil
}
