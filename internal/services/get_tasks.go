package services

import (
	"log/slog"
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
	slog.Info("get tasks service started...")
	user, _ := s.userRepository.GetFirstById(userId)

	if user == nil {
		slog.Error("user not found error")
		return []entities.Task{}
	}

	tasks := s.taskRepository.GetByUserId(user.GetId())

	slog.Info("get tasks service finished...")

	return tasks
}
