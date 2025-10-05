package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/interfaces"
	"scheduler/internal/queue"
	"scheduler/pkg/scheduler"
	"time"
)

type CreateTaskService struct {
	userRepository           interfaces.IUserRepository
	taskRepository           interfaces.ITaskRepository
	createTransactionService *CreateTransactionService
	updateTransactionService *UpdateTaskTransactionService
	scheduler                *scheduler.Scheduler
}

func NewCreateTaskService(
	userRepository interfaces.IUserRepository,
	taskRepository interfaces.ITaskRepository,
	createTransactionService *CreateTransactionService,
	updateTransactionService *UpdateTaskTransactionService,
	scheduler *scheduler.Scheduler,
) *CreateTaskService {
	return &CreateTaskService{
		userRepository:           userRepository,
		taskRepository:           taskRepository,
		createTransactionService: createTransactionService,
		updateTransactionService: updateTransactionService,
		scheduler:                scheduler,
	}
}

func (s *CreateTaskService) Execute(
	kind string,
	runAt time.Time,
	timezone string,
	priority int,
	userId string,
	referenceId string,
	idempotencyKey string,
) (*entities.Task, error) {
	slog.Info("create task service started...")
	slog.Debug(fmt.Sprint("input ", kind,
		runAt,
		timezone,
		priority,
		userId,
		referenceId,
		idempotencyKey))

	if _, ok := queue.WorkersQueues[kind]; !ok {
		reason := "the type is not valid"
		return nil, errors.INVALID_FIELD_VALUE("type", &reason)
	}

	user, _ := s.userRepository.GetFirstById(userId)
	if user == nil {
		slog.Error("user not found error")
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	task, _ := s.taskRepository.GetFirstByIdempotencyKey(idempotencyKey)

	if task != nil {
		slog.Info("task already exists with this idempotency key")
		return task, nil
	}

	task, err := entities.NewTask(
		kind,
		user.GetId(),
		10,
		runAt,
		timezone,
		referenceId,
		idempotencyKey,
	)
	if err != nil {
		slog.Error(fmt.Sprintf("task creation error %s", err.Error()))
		return nil, err
	}

	if priority != entities.PriorityLow {
		err = task.SetPriority(priority)
		if err != nil {
			slog.Error(fmt.Sprintf("priority set error %s", err.Error()))
			return nil, err
		}
	}

	err = s.taskRepository.Create(task)
	if err != nil {
		slog.Error(fmt.Sprintf("save task into repository error %s", err.Error()))
		return nil, err
	}

	transaction, err := s.createTransactionService.Execute(
		user.GetId(),
		task.GetCost(),
		"TASK",
		entities.TypeTransactionTaskSend,
		task.GetId(),
		task.GetIdempotencyKey(),
	)
	if err != nil {
		slog.Error(fmt.Sprintf("transaction creation error %s", err.Error()))
		return nil, err
	}

	_, err = s.updateTransactionService.Frozen(transaction.GetId())
	if err != nil {
		slog.Error(fmt.Sprintf("update transaction error %s", err.Error()))
		return nil, err
	}

	s.scheduler.Add(task)

	slog.Info("create task service finished...")
	slog.Debug(fmt.Sprintf("returned task %+v", task))

	return task, nil
}
