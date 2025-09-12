package scheduler

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/interfaces"
	"scheduler/internal/persistence"

	"github.com/jonboulle/clockwork"
)

type Scheduler struct {
	heap     Heap
	capacity int
	clock    clockwork.Clock
	addCh    chan *entities.Task
	stopCh   chan struct{}
	queue    interfaces.IQueue
}

var indexes = make(map[string]int)

func NewScheduler(clock clockwork.Clock, queue interfaces.IQueue, capacity int) *Scheduler {
	sc := &Scheduler{
		heap:     Heap{},
		capacity: capacity,
		clock:    clock,
		addCh:    make(chan *entities.Task),
		stopCh:   make(chan struct{}),
		queue:    queue,
	}

	heap.Init(&sc.heap)
	return sc
}

func (s *Scheduler) PeekTasks() Heap {
	return s.heap
}

func (s *Scheduler) Run() {
	for {
		if len(s.heap) == 0 {
			select {
			case x := <-s.addCh:
				s.upsert(x)
			case <-s.stopCh:
				return
			}
		} else {
			next := s.heap[0]
			now := s.clock.Now()
			waitDuration := next.GetRunAt().Sub(now)

			if waitDuration <= 0 {
				x := heap.Pop(&s.heap).(*entities.Task)

				s.ExecuteTask(x)
			} else {
				timer := s.clock.After(waitDuration)

				select {
				case <-timer:
					task := heap.Pop(&s.heap).(*entities.Task)
					s.ExecuteTask(task)
				case task := <-s.addCh:
					s.upsert(task)
				case <-s.stopCh:
					return
				}
			}
		}
	}
}

func (s *Scheduler) ExecuteTask(task *entities.Task) {
	slog.Info(
		fmt.Sprintf(
			"task %s sent to queue %s at %s",
			task.GetId(),
			task.GetType(),
			task.GetRunAt().String(),
		),
	)

	data, _ := json.Marshal(persistence.ToTaskModel(task))

	_ = s.queue.Publish("tasks", data)
}

func (s *Scheduler) upsert(task *entities.Task) {
	index, ok := indexes[task.GetId()]

	slog.Debug(fmt.Sprintf("current indexes %v", indexes))

	if ok {
		s.heap[index] = task
		heap.Fix(&s.heap, index)

		slog.Info(fmt.Sprintf("task updated into index %d, task id %s", index, task.GetId()))
		return
	}

	heap.Push(&s.heap, task)
	indexes[task.GetId()] = len(s.heap) - 1
	slog.Info(fmt.Sprintf("new task %s added to heap", task.GetId()))
}

func (s *Scheduler) Add(task *entities.Task) {
	s.addCh <- task
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}
