package services

import (
	"scheduler/internal/entities"
	repos "scheduler/internal/repositories"
)

type GetTasksService struct {
	userRepository repos.IUserRepository
	taskRepository repos.ITaskRepository
}

func NewGetTasksService(
	userRepository repos.IUserRepository,
	taskRepository repos.ITaskRepository,
) *GetTasksService {
	return &GetTasksService{
		userRepository: userRepository,
		taskRepository: taskRepository,
	}
}

func (s *GetTasksService) Execute(userId string) []entities.Task {
	user, _ := s.userRepository.GetFirstById(userId)

	if user == nil {
		return []entities.Task{}
	}

	tasks := s.taskRepository.GetByUserId(user.GetId())

	return tasks
}
