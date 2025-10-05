package queue_handlers

import (
	"context"
	"encoding/json"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/interfaces"
	"scheduler/internal/queue"
	"scheduler/internal/services"
	"scheduler/pkg/postgres"
)

type UpdateTaskHandler struct {
	db                *postgres.Database
	queue             interfaces.IQueue
	updateTaskService *services.UpdateTaskService
}

func NewUpdateTaskHandler(
	db *postgres.Database,
	queue interfaces.IQueue,
	updateTaskService *services.UpdateTaskService,
) *UpdateTaskHandler {
	return &UpdateTaskHandler{
		db:                db,
		queue:             queue,
		updateTaskService: updateTaskService,
	}
}

func (h *UpdateTaskHandler) Handle(ctx context.Context) error {
	return h.queue.Consume(ctx, queue.GET_TASKS_RESULT_QUEUE, func(msg any) error {
		data, ok := msg.(queue.TaskUpdate)

		if !ok {
			return errors.INVALID_FIELD_VALUE("msg", nil)
		}

		switch data.Status {
		case entities.StatusCompleted:
			h.db.BeginTransaction()

			task, err := h.updateTaskService.Complete(data.ID)
			if err != nil {
				_ = h.db.RollbackTransaction()
				return err
			}

			_ = h.db.CommitTransaction()

			event, err := queue.NewEvent(task.GetUserId(), "success", map[string]any{
				"taskId":   task.GetId(),
				"when":     task.GetRunAt(),
				"timezone": task.GetTimezone(),
				"type":     task.GetType(),
			})
			if err != nil {
				return err
			}

			data, err := json.Marshal(event)
			if err != nil {
				return err
			}

			err = h.queue.Publish(queue.SEND_EVENTS_KEY, queue.EVENTS_EXCHANGE, data)
			if err != nil {
				return err
			}

		case entities.StatusFailed:
			h.db.BeginTransaction()
			task, err := h.updateTaskService.Fail(data.ID, *data.Refund, *data.Reason)
			if err != nil {
				_ = h.db.RollbackTransaction()
				return err
			}
			_ = h.db.CommitTransaction()

			event, err := queue.NewEvent(task.GetUserId(), "error", map[string]any{
				"message":  *data.Reason,
				"refund":   *data.Refund,
				"when":     task.GetRunAt(),
				"timezone": task.GetTimezone(),
			})
			if err != nil {
				return err
			}

			data, err := json.Marshal(event)
			if err != nil {
				return err
			}

			err = h.queue.Publish(queue.SEND_EVENTS_KEY, queue.EVENTS_EXCHANGE, data)
			if err != nil {
				return err
			}

		}

		return nil
	})
}
