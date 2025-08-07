package entities_test

import (
	"scheduler/internal/entities"
	. "scheduler/test"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTask_New(t *testing.T) {
	task, err := entities.NewTask("email_send", uuid.NewString(), time.Now().Add(1*time.Hour))

	Ok(t, err)
	Equals(t, entities.TaskCost["email_send"], task.GetCost())
	Equals(t, "pending", task.GetStatus())
}

func TestTask_WrongTaskType(t *testing.T) {
	_, err := entities.NewTask("email_sent", uuid.NewString(), time.Now().Add(1*time.Hour))

	Assert(t, err != nil, "expect error for invalid type got success")
}

func TestTask_SetStatus(t *testing.T) {
	task, err := entities.NewTask("email_send", uuid.NewString(), time.Now().Add(1*time.Hour))

	Ok(t, err)
	Equals(t, "pending", task.GetStatus())

	err = task.SetStatus("completed")

	Ok(t, err)
}

func TestTask_SetStatusAfterCompleted(t *testing.T) {
	task, err := entities.NewTask("email_send", uuid.NewString(), time.Now().Add(1*time.Hour))

	Ok(t, err)
	Equals(t, "pending", task.GetStatus())

	err = task.SetStatus("completed")

	Ok(t, err)

	err = task.SetStatus("pending")

	Assert(t, err != nil, "expect error in set status from complete to pending got success")
}

func TestTask_SetStatusAfterFailed(t *testing.T) {
	task, err := entities.NewTask("email_send", uuid.NewString(), time.Now().Add(1*time.Hour))

	Ok(t, err)
	Equals(t, "pending", task.GetStatus())

	err = task.SetStatus("failed")

	Ok(t, err)

	err = task.SetStatus("pending")

	Assert(t, err != nil, "expect error in set status from failed to pending got success")
}

func TestTask_InvalidUserId(t *testing.T) {
	task, err := entities.NewTask("email_send", uuid.NewString(), time.Now().Add(1*time.Hour))

	Ok(t, err)

	err = task.SetUserId("42142141")

	Assert(t, err != nil, "expect invalid user id error got success")
}

func TestTask_TimeAfterOneYear(t *testing.T) {
	task, err := entities.NewTask("email_send", uuid.NewString(), time.Now().Add(1*time.Hour))

	Ok(t, err)

	err = task.SetTime(time.Now().AddDate(2, 0, 1))

	Assert(t, err != nil, "expect error cause the time after one year got success")
}

func TestTask_TimeBeforeNow(t *testing.T) {
	task, err := entities.NewTask("email_send", uuid.NewString(), time.Now().Add(1*time.Hour))

	Ok(t, err)

	err = task.SetTime(time.Now().AddDate(0, 0, -1))

	Assert(t, err != nil, "expect error cause the time before now got success")
}
