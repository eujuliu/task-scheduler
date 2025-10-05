package workers

import (
	"encoding/json"
	"scheduler/internal/entities"
	"scheduler/internal/interfaces"
	q "scheduler/internal/queue"
	"time"
)

type EmailWorkerRequest struct {
	ID          string    `json:"id"`
	UserId      string    `json:"user_id"`
	RunAt       time.Time `json:"run_at"`
	ReferenceId string    `json:"reference_id"`
}

func EmailWorker(queue interfaces.IQueue, task *entities.Task) {
	req := EmailWorkerRequest{
		ID:          task.GetId(),
		UserId:      task.GetUserId(),
		RunAt:       task.GetRunAt(),
		ReferenceId: task.GetReferenceId(),
	}

	data, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	err = queue.Publish(q.WorkersQueues["email"], q.TASK_EXCHANGE, data, req.ID)
	if err != nil {
		panic(err)
	}
}
