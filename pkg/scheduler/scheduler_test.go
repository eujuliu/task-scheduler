package scheduler_test

import (
	"encoding/json"
	"scheduler/internal/entities"
	"scheduler/internal/persistence"
	"scheduler/internal/queue"
	"scheduler/pkg/scheduler"
	. "scheduler/test"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
)

func TestScheduler(t *testing.T) {
	queue := queue.NewInMemoryQueue()
	c := clockwork.NewFakeClock()

	sc := scheduler.NewScheduler(c, queue, 20)

	go sc.Run()

	times := make([]time.Time, 5)
	times[0] = c.Now().Add(5 * time.Minute)

	for i := 1; i < 5; i++ {
		times[i] = times[i-1].Add(1 * time.Minute)
	}

	task, _ := entities.NewTask(
		"video",
		uuid.NewString(),
		10,
		times[0],
		"America/Sao_Paulo",
		"0",
		"123",
	)
	sc.Add(task)

	task1, _ := entities.NewTask(
		"video",
		uuid.NewString(),
		10,
		times[4],
		"America/Sao_Paulo",
		"1",
		"123",
	)
	sc.Add(task1)

	task2, _ := entities.NewTask(
		"video",
		uuid.NewString(),
		10,
		times[3],
		"America/Sao_Paulo",
		"2",
		"123",
	)
	sc.Add(task2)

	task3, _ := entities.NewTask(
		"video",
		uuid.NewString(),
		10,
		times[2],
		"America/Sao_Paulo",
		"3",
		"123",
	)
	sc.Add(task3)

	task4, _ := entities.NewTask(
		"video",
		uuid.NewString(),
		10,
		times[1],
		"America/Sao_Paulo",
		"4",
		"123",
	)
	sc.Add(task4)

	tasksCh, _ := queue.Consume("tasks")

	select {
	case <-tasksCh:
		t.Fatal("unexpected received result before timer expired")
	case <-time.After(2 * time.Second):
	}

	c.Advance(5 * time.Minute)

	result := []string{"0", "4", "3", "2", "1"}

	for _, i := range result {
		d := <-tasksCh
		var got persistence.TaskModel
		_ = json.Unmarshal(d, &got)

		Equals(t, i, got.ReferenceID)
		Equals(t, c.Now().Format("23:00:00"), got.RunAt.Format("23:00:00"))
		c.Advance(1 * time.Minute)
	}
}
