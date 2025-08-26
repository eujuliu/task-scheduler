package services

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"time"

	repos "scheduler/internal/repositories"
)

type CreateTaskService struct {
	userRepository           repos.IUserRepository
	transactionRepository    repos.ITransactionRepository
	taskRepository           repos.ITaskRepository
	createTransactionService *CreateTransactionService
	updateTransactionService *UpdateTaskTransactionService
}

func NewCreateTaskService(
	userRepository repos.IUserRepository,
	transactionRepository repos.ITransactionRepository,
	taskRepository repos.ITaskRepository,
	createTransactionService *CreateTransactionService,
	updateTransactionService *UpdateTaskTransactionService,
) *CreateTaskService {
	return &CreateTaskService{
		userRepository:           userRepository,
		transactionRepository:    transactionRepository,
		taskRepository:           taskRepository,
		createTransactionService: createTransactionService,
		updateTransactionService: updateTransactionService,
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

	task, err = entities.NewTask(
		kind,
		user.GetId(),
		10,
		runAt,
		timezone,
		referenceId,
		idempotencyKey,
	)
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

	transaction, err := s.createTransactionService.Execute(
		user.GetId(),
		task.GetCost(),
		0,
		"",
		entities.TypeTransactionTaskSend,
		task.GetReferenceId(),
		task.GetIdempotencyKey(),
	)
	if err != nil {
		return nil, err
	}

	_, err = s.updateTransactionService.Frozen(transaction.GetId())
	if err != nil {
		return nil, err
	}

	return task, nil
}
