package services

import (
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
	user, _ := s.userRepository.GetFirstById(userId)

	if user == nil {
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	task, _ := s.taskRepository.GetFirstById(taskId)

	if task == nil {
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	return task, nil
}
