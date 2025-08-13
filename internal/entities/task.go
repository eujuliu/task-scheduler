package entities

import (
	"scheduler/internal/errors"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	BaseEntity
	kind        string
	status      string
	userId      string
	time        time.Time
	cost        int
	referenceId string
}

var TaskCost = map[string]int{
	"email_send": 10,
}

func NewTask(kind string, userId string, referenceId string, scheduledTime time.Time) (*Task, error) {
	taskType := strings.TrimSpace(strings.ToLower(kind))
	_, ok := TaskCost[taskType]

	if !ok {
		return nil, errors.INVALID_FIELD_VALUE("task type")
	}

	task := &Task{
		BaseEntity:  *NewBaseEntity(),
		kind:        taskType,
		status:      "pending",
		cost:        TaskCost[taskType],
		referenceId: referenceId,
	}

	err := task.SetTime(scheduledTime)

	if err != nil {
		return nil, err
	}

	err = task.SetUserId(userId)

	if err != nil {
		return nil, err
	}

	return task, nil

}

func (t *Task) GetType() string {
	return t.kind
}

func (t *Task) SetStatus(status string) error {
	if t.readonly() {
		return errors.FINISHED_OPERATION_ERROR()
	}

	formatted := strings.TrimSpace(strings.ToLower(status))

	var AvailableStatus = []string{"pending", "running", "completed", "failed"}

	if !slices.Contains(AvailableStatus, formatted) {
		return errors.INVALID_FIELD_VALUE("status")
	}

	t.status = formatted

	return nil
}

func (t *Task) GetStatus() string {
	return t.status
}

func (t *Task) SetUserId(id string) error {
	if t.readonly() {
		return errors.FINISHED_OPERATION_ERROR()
	}

	if err := uuid.Validate(id); err != nil {
		return errors.INVALID_FIELD_VALUE("user id")
	}

	t.userId = id

	return nil
}

func (t *Task) GetUserId() string {
	return t.userId
}

func (t *Task) SetTime(when time.Time) error {
	if t.readonly() {
		return errors.FINISHED_OPERATION_ERROR()
	}

	now := time.Now()
	afterOneYear := now.AddDate(1, 0, 0)

	if when.Before(now) || when.After(afterOneYear) {
		return errors.INVALID_FIELD_VALUE("schedule time")
	}

	t.time = when

	return nil
}

func (t *Task) GetTime() time.Time {
	return t.time
}

func (t *Task) GetCost() int {
	return t.cost
}

func (t *Task) GetReferenceId() string {
	return t.referenceId
}

func (t *Task) readonly() bool {
	unchangeableStatus := []string{"completed", "failed"}

	return slices.Contains(unchangeableStatus, t.status)
}
