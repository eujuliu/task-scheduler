package scheduler_test

import (
	"container/heap"
	"fmt"
	"scheduler/internal/entities"
	"scheduler/pkg/scheduler"
	"testing"
	"time"

	"github.com/google/uuid"

	. "scheduler/test"
)

func TestHeap_Ordering(t *testing.T) {
	h := &scheduler.Heap{}

	now := time.Now()

	heap.Init(h)

	for i := range 5 {
		t, _ := entities.NewTask(
			"video",
			uuid.NewString(),
			10,
			now.Add(time.Duration(i+2)*time.Second),
			"America/Sao_Paulo",
			fmt.Sprint(5-i),
			"123",
		)

		heap.Push(h, t)
	}

	expected := []string{"5", "4", "3", "2", "1"}

	for i, got := range *h {
		want := expected[i]

		Equals(t, want, got.GetReferenceId())
	}
}
