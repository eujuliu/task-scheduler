package scheduler_test

import (
	"context"
	"scheduler/internal/entities"
	"scheduler/internal/queue"
	. "scheduler/test"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestScheduler(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	ctx := context.Background()

	times := make([]time.Time, 5)
	times[0] = Clock.Now().Add(5 * time.Minute)

	for i := 1; i < 5; i++ {
		times[i] = times[i-1].Add(1 * time.Minute)
	}

	task, _ := entities.NewTask(
		"email",
		uuid.NewString(),
		10,
		times[0],
		"America/Sao_Paulo",
		"0",
		"123",
	)
	Scheduler.Add(task)

	task1, _ := entities.NewTask(
		"email",
		uuid.NewString(),
		10,
		times[4],
		"America/Sao_Paulo",
		"1",
		"123",
	)
	Scheduler.Add(task1)

	task2, _ := entities.NewTask(
		"email",
		uuid.NewString(),
		10,
		times[3],
		"America/Sao_Paulo",
		"2",
		"123",
	)
	Scheduler.Add(task2)

	task3, _ := entities.NewTask(
		"email",
		uuid.NewString(),
		10,
		times[2],
		"America/Sao_Paulo",
		"3",
		"123",
	)
	Scheduler.Add(task3)

	task4, _ := entities.NewTask(
		"email",
		uuid.NewString(),
		10,
		times[1],
		"America/Sao_Paulo",
		"4",
		"123",
	)
	Scheduler.Add(task4)

	tasksCh := make(chan any)

	go func() {
		err := Queue.Consume(ctx, queue.SEND_EMAIL_KEY, func(message any) error {
			tasksCh <- message

			return nil
		})

		Ok(t, err)
	}()

	select {
	case <-tasksCh:
		t.Fatal("unexpected received result before timer expired")
	case <-time.After(2 * time.Second):
	}

	Clock.Advance(5 * time.Minute)

	result := []string{"0", "4", "3", "2", "1"}

	for _, i := range result {
		msg := <-tasksCh
		got := msg.(map[string]any)
		date, err := time.Parse(time.RFC3339Nano, got["runAt"].(string))

		Ok(t, err)
		Equals(t, i, got["referenceId"])
		Equals(t, Clock.Now().Format("23:00:00"), date.Format("23:00:00"))
		Equals(t, entities.StatusRunning, got["status"])

		Clock.Advance(1 * time.Minute)
	}
}
