package services

import (
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/interfaces"
)

type GetTasksByUserIdService struct {
	userRepository interfaces.IUserRepository
	taskRepository interfaces.ITaskRepository
}

func NewGetTasksByUserIdService(
	userRepository interfaces.IUserRepository,
	taskRepository interfaces.ITaskRepository,
) *GetTasksByUserIdService {
	return &GetTasksByUserIdService{
		userRepository: userRepository,
		taskRepository: taskRepository,
	}
}

func (s *GetTasksByUserIdService) Execute(
	userId string,
	offset *int,
	limit *int,
	orderBy *string,
) []entities.Task {
	slog.Info("get tasks by user id service started...")
	user, _ := s.userRepository.GetFirstById(userId)

	if user == nil {
		slog.Error("user not found error")
		return []entities.Task{}
	}

	tasks := s.taskRepository.GetByUserId(user.GetId(), offset, limit, orderBy)

	slog.Info("get tasks by user id service finished...")

	return tasks
}
