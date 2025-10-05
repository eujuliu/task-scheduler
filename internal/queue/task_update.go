package queue

type TaskUpdate struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Reason *string `json:"reason"`
	Refund *bool   `json:"refund"`
}

func NewTaskUpdate(id, status string, reason *string, refund *bool) *TaskUpdate {
	res := &TaskUpdate{
		ID:     id,
		Status: status,
		Reason: reason,
		Refund: refund,
	}

	return res
}
