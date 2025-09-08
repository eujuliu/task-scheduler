package entities

import (
	"fmt"
	"scheduler/internal/errors"
	"slices"
	"time"
)

const (
	PriorityHigh = iota + 1
	PriorityMedium
	PriorityLow
	PriorityExtraLow
)

type Task struct {
	BaseEntity
	kind           string
	userId         string
	cost           int
	status         string
	runAt          time.Time
	timezone       string
	retries        int
	priority       int
	referenceId    string
	idempotencyKey string
}

func NewTask(
	kind string,
	userId string,
	cost int,
	runAt time.Time,
	timezone string,
	referenceId string,
	idempotencyKey string,
) (*Task, error) {
	task := &Task{
		BaseEntity:     *NewBaseEntity(),
		kind:           kind,
		userId:         userId,
		cost:           cost,
		status:         StatusPending,
		retries:        0,
		priority:       PriorityLow,
		referenceId:    referenceId,
		idempotencyKey: idempotencyKey,
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

func HydrateTask(
	id, kind, userId string,
	cost int,
	status string,
	runAt time.Time,
	timezone string,
	retries, priority int,
	referenceId, idempotencyKey string,
	createdAt, updateAt time.Time,
) *Task {
	return &Task{
		BaseEntity: BaseEntity{
			id:        id,
			createdAt: createdAt,
			updatedAt: updateAt,
		},

		kind:           kind,
		userId:         userId,
		cost:           cost,
		status:         status,
		runAt:          runAt,
		timezone:       timezone,
		retries:        retries,
		priority:       priority,
		referenceId:    referenceId,
		idempotencyKey: idempotencyKey,
	}
}

func (t *Task) GetType() string {
	return t.kind
}

func (t *Task) GetUserId() string {
	return t.userId
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
		reason := fmt.Sprintf(
			"you can set one of these (%s, %s, %s, %s)",
			StatusRunning,
			StatusCompleted,
			StatusFailed,
			StatusFrozen,
		)
		return errors.INVALID_FIELD_VALUE("status", &reason)
	}

	t.status = status

	return nil
}

func (t *Task) GetStatus() string {
	return t.status
}

func (t *Task) SetRunAt(when time.Time) error {
	if t.readonly() || t.running() {
		return errors.FINISHED_OPERATION_ERROR()
	}

	now := time.Now().UTC()

	if when.Before(now) || when.After(now.AddDate(0, 6, 0)) {
		return errors.INVALID_FIELD_VALUE("run at", nil)
	}

	t.runAt = when

	return nil
}

func (t *Task) GetRunAt() time.Time {
	return t.runAt
}

func (t *Task) SetTimezone(timezone string) error {
	if t.readonly() || t.running() {
		return errors.FINISHED_OPERATION_ERROR()
	}

	_, err := time.LoadLocation(timezone)
	if err != nil {
		reason := "you need to set one of the IATA timezones"
		return errors.INVALID_FIELD_VALUE("timezone", &reason)
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
	if t.readonly() || t.running() {
		return errors.FINISHED_OPERATION_ERROR()
	}

	availablePriorities := []int{
		PriorityHigh,
		PriorityMedium,
		PriorityLow,
		PriorityExtraLow,
	}

	if !slices.Contains(availablePriorities, priority) {
		reason := fmt.Sprintf("you need to set one of these (%d, %d, %d, %d)", PriorityHigh,
			PriorityMedium,
			PriorityLow,
			PriorityExtraLow)
		return errors.INVALID_FIELD_VALUE("priority", &reason)
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

func (t *Task) GetIdempotencyKey() string {
	return t.idempotencyKey
}

func (t *Task) running() bool {
	return t.status == StatusRunning
}

func (t *Task) readonly() bool {
	return slices.Contains(
		[]string{StatusCompleted, StatusFailed, StatusCanceled},
		t.status,
	)
}
