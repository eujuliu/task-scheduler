package entities

import (
	"scheduler/internal/errors"
	"slices"
	"time"
)

const (
	StatusPending   string = "PENDING"
	StatusRunning   string = "RUNNING"
	StatusCompleted string = "COMPLETED"
	StatusFailed    string = "FAILED"
	StatusCanceled  string = "CANCELED"
)

const (
	PriorityHigh = iota + 1
	PriorityMedium
	PriorityLow
	PriorityExtraLow
)

type Task struct {
	BaseEntity
	kind        string
	userId      string
	cost        int
	status      string
	runAt       time.Time
	timezone    string
	retries     int
	priority    int
	referenceId string
}

func NewTask(kind string, userId string, cost int, runAt time.Time, timezone string, referenceId string) (*Task, error) {
	task := &Task{
		BaseEntity:  *NewBaseEntity(),
		kind:        kind,
		userId:      userId,
		cost:        cost,
		status:      StatusPending,
		retries:     0,
		priority:    PriorityLow,
		referenceId: referenceId,
	}

	err := task.SetRunAt(runAt)

	if err != nil {
		return nil, err
	}

	err = task.SetTimezone(timezone)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *Task) GetType() string {
	return t.kind
}

func (t *Task) GetCost() int {
	return t.cost
}

func (t *Task) SetStatus(status string) error {
	if t.readonly() {
		return errors.FINISHED_OPERATION_ERROR()
	}

	availableStatus := []string{
		StatusRunning,
		StatusCompleted,
		StatusFailed,
		StatusCanceled,
	}

	if !slices.Contains(availableStatus, status) {
		return errors.INVALID_FIELD_VALUE("status")
	}

	t.status = status

	return nil
}

func (t *Task) GetStatus() string {
	return t.status
}

func (t *Task) SetRunAt(when time.Time) error {
	if t.readonly() {
		return errors.FINISHED_OPERATION_ERROR()
	}

	now := time.Now().UTC()

	if when.Before(now) || when.After(now.AddDate(0, 6, 0)) {
		return errors.INVALID_FIELD_VALUE("run at")
	}

	t.runAt = when

	return nil
}

func (t *Task) GetRunAt() time.Time {
	return t.runAt
}

func (t *Task) SetTimezone(timezone string) error {
	if t.readonly() {
		return errors.FINISHED_OPERATION_ERROR()
	}

	_, err := time.LoadLocation(timezone)

	if err != nil {
		return errors.INVALID_FIELD_VALUE("timezone")
	}

	t.timezone = timezone

	return nil
}

func (t *Task) GetTimezone() string {
	return t.timezone
}

func (t *Task) AddRetry() error {
	if t.readonly() {
		return errors.FINISHED_OPERATION_ERROR()
	}

	if t.retries == 3 {
		return errors.MAX_VALUE_REACHED_ERROR()
	}

	t.retries += 1

	return nil
}

func (t *Task) GetRetries() int {
	return t.retries
}

func (t *Task) SetPriority(priority int) error {
	if t.readonly() {
		return errors.FINISHED_OPERATION_ERROR()
	}

	availablePriorities := []int{
		PriorityHigh,
		PriorityMedium,
		PriorityLow,
		PriorityExtraLow,
	}

	if !slices.Contains(availablePriorities, priority) {
		return errors.INVALID_FIELD_VALUE("priority")
	}

	t.priority = priority

	return nil
}

func (t *Task) GetPriority() int {
	return t.priority
}

func (t *Task) GetReferenceId() string {
	return t.referenceId
}

func (t *Task) readonly() bool {
	return slices.Contains([]string{StatusCompleted, StatusFailed, StatusCanceled}, t.status)
}
