package services

import (
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/interfaces"
	"time"
)

type GetTasksByRunAtService struct {
	taskRepository interfaces.ITaskRepository
}

func NewGetTasksByRunAtService(
	taskRepository interfaces.ITaskRepository,
) *GetTasksByRunAtService {
	return &GetTasksByRunAtService{
		taskRepository: taskRepository,
	}
}

func (s *GetTasksByRunAtService) Execute(
	status string,
	asc bool,
	limit int,
	from *time.Time,
) []entities.Task {
	slog.Info("get tasks by run at service started...")
	slog.Debug(fmt.Sprint("input, ", status, asc, limit, from))

	tasks := s.taskRepository.Get(&status, &asc, &limit, from)

	slog.Info("get tasks by run at service finished...")

	return tasks
}
