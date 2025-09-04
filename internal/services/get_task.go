package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	repos "scheduler/internal/repositories"
)

type GetTaskService struct {
	userRepository repos.IUserRepository
	taskRepository repos.ITaskRepository
}

func NewGetTaskService(
	userRepository repos.IUserRepository,
	taskRepository repos.ITaskRepository,
) *GetTaskService {
	return &GetTaskService{
		userRepository: userRepository,
		taskRepository: taskRepository,
	}
}

func (s *GetTaskService) Execute(userId string, taskId string) (*entities.Task, error) {
	slog.Info("get task service started...")
	slog.Debug(fmt.Sprint("input ", userId, taskId))
	user, _ := s.userRepository.GetFirstById(userId)

	if user == nil {
		slog.Error("user not found error")
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	task, _ := s.taskRepository.GetFirstById(taskId)

	if task == nil {
		slog.Error("task not found error")
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	slog.Info("get task service finished...")
	slog.Debug(fmt.Sprintf("returned task %+v", task))

	return task, nil
}
