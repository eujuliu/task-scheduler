package entities_test

import (
	. "scheduler/internal/entities"
	"scheduler/internal/errors"
	. "scheduler/test"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTask(t *testing.T) {
	task, err := NewTask("email", uuid.NewString(), 20, time.Now().UTC().AddDate(0, 0, 2), "America/Sao_Paulo", uuid.NewString())

	Ok(t, err)
	Equals(t, 3, task.GetPriority())
	Equals(t, 0, task.GetRetries())
	Equals(t, StatusPending, task.GetStatus())
}

func TestTask_SetStatus(t *testing.T) {
	task, err := NewTask("email", uuid.NewString(), 20, time.Now().UTC().AddDate(0, 0, 2), "America/Sao_Paulo", uuid.NewString())

	Ok(t, err)

	err = task.SetStatus(StatusRunning)

	Ok(t, err)
}

func TestTask_SetStatusAfterCompleted(t *testing.T) {
	task, err := NewTask("email", uuid.NewString(), 20, time.Now().UTC().AddDate(0, 0, 2), "America/Sao_Paulo", uuid.NewString())

	Ok(t, err)

	err = task.SetStatus(StatusCompleted)

	Ok(t, err)

	err = task.SetStatus(StatusCanceled)

	Assert(t, err != nil, "expected error got success")
	Equals(t, errors.FINISHED_OPERATION_ERROR().Error(), err.Error())
}

func TestTask_SetRunAt(t *testing.T) {
	task, err := NewTask("email", uuid.NewString(), 20, time.Now().UTC().AddDate(0, 0, 2), "America/Sao_Paulo", uuid.NewString())

	Ok(t, err)

	err = task.SetRunAt(time.Now().UTC().AddDate(0, 0, 3))

	Ok(t, err)
}

func TestTask_SetRunAtInCompletedTask(t *testing.T) {
	task, err := NewTask("email", uuid.NewString(), 20, time.Now().UTC().AddDate(0, 0, 2), "America/Sao_Paulo", uuid.NewString())

	Ok(t, err)

	err = task.SetStatus(StatusCompleted)

	Ok(t, err)

	err = task.SetRunAt(time.Now().UTC().AddDate(0, 0, 3))

	Assert(t, err != nil, "expected err got success")
	Equals(t, errors.FINISHED_OPERATION_ERROR().Error(), err.Error())
}

func TestTask_SetRunAtBeforeToday(t *testing.T) {
	_, err := NewTask("email", uuid.NewString(), 20, time.Now().UTC().AddDate(0, 0, -2), "America/Sao_Paulo", uuid.NewString())

	Assert(t, err != nil, "expected err got success")
	Equals(t, errors.INVALID_FIELD_VALUE("run at").Error(), err.Error())
}

func TestTask_SetRunAtAfterSixMonthsForward(t *testing.T) {
	_, err := NewTask("email", uuid.NewString(), 20, time.Now().UTC().AddDate(0, 7, 0), "America/Sao_Paulo", uuid.NewString())

	Assert(t, err != nil, "expected err got success")
	Equals(t, errors.INVALID_FIELD_VALUE("run at").Error(), err.Error())
}

func TestTask_SetTimezone(t *testing.T) {
	task, err := NewTask("email", uuid.NewString(), 20, time.Now().UTC().AddDate(0, 0, 2), "America/Sao_Paulo", uuid.NewString())

	Ok(t, err)

	err = task.SetTimezone("Europe/Berlin")

	Ok(t, err)
}

func TestTask_SetTimezoneInvalidTimezone(t *testing.T) {
	task, err := NewTask("email", uuid.NewString(), 20, time.Now().UTC().AddDate(0, 0, 2), "America/Sao_Paulo", uuid.NewString())

	Ok(t, err)

	err = task.SetTimezone("GTM-3")

	Assert(t, err != nil, "expected error got success")
	Equals(t, errors.INVALID_FIELD_VALUE("timezone").Error(), err.Error())
}

func TestTask_AddRetry(t *testing.T) {
	task, err := NewTask("email", uuid.NewString(), 20, time.Now().UTC().AddDate(0, 0, 2), "America/Sao_Paulo", uuid.NewString())

	Ok(t, err)
	Equals(t, 0, task.GetRetries())

	err = task.AddRetry()

	Ok(t, err)
	Equals(t, 1, task.GetRetries())

	err = task.AddRetry()

	Ok(t, err)
	Equals(t, 2, task.GetRetries())
}

func TestTask_AddRetryReachMax(t *testing.T) {
	task, err := NewTask("email", uuid.NewString(), 20, time.Now().UTC().AddDate(0, 0, 2), "America/Sao_Paulo", uuid.NewString())

	Ok(t, err)
	Equals(t, 0, task.GetRetries())

	err = task.AddRetry()

	Ok(t, err)
	Equals(t, 1, task.GetRetries())

	err = task.AddRetry()

	Ok(t, err)
	Equals(t, 2, task.GetRetries())

	err = task.AddRetry()

	Ok(t, err)
	Equals(t, 3, task.GetRetries())

	err = task.AddRetry()

	Assert(t, err != nil, "expected error got success")
	Equals(t, errors.MAX_VALUE_REACHED_ERROR().Error(), err.Error())
}

func TestTask_SetPriority(t *testing.T) {
	task, err := NewTask("email", uuid.NewString(), 20, time.Now().UTC().AddDate(0, 0, 2), "America/Sao_Paulo", uuid.NewString())

	Ok(t, err)
	Equals(t, PriorityLow, task.GetPriority())

	err = task.SetPriority(PriorityHigh)

	Ok(t, err)
	Equals(t, PriorityHigh, task.GetPriority())
}
