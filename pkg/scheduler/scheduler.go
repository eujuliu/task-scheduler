package scheduler

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log/slog"
	"scheduler/internal/entities"
	"scheduler/internal/interfaces"
	"scheduler/internal/persistence"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
)

type Scheduler struct {
	heap           Heap
	capacity       int
	mu             sync.Mutex
	clock          clockwork.Clock
	addCh          chan *entities.Task
	stopCh         chan struct{}
	queue          interfaces.IQueue
	taskRepository interfaces.ITaskRepository
}

var indexes = make(map[string]int)

func NewScheduler(
	clock clockwork.Clock,
	queue interfaces.IQueue,
	capacity int,
	taskRepository interfaces.ITaskRepository,
) *Scheduler {
	sc := &Scheduler{
		heap:           Heap{},
		capacity:       capacity,
		clock:          clock,
		addCh:          make(chan *entities.Task),
		stopCh:         make(chan struct{}),
		queue:          queue,
		taskRepository: taskRepository,
	}

	heap.Init(&sc.heap)

	return sc
}

func (s *Scheduler) PeekTasks() Heap {
	return s.heap
}

func (s *Scheduler) Run() {
	go s.loadNextTasks()

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

	err := task.SetStatus(entities.StatusRunning)
	if err != nil {
		panic(err)
	}

	err = s.taskRepository.Update(task)
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(persistence.ToTaskModel(task))
	if err != nil {
		panic(err)
	}

	err = s.queue.Publish("task.send", "task-exchange", data)
	if err != nil {
		panic(err)
	}

	go s.loadNextTasks()
}

func (s *Scheduler) loadNextTasks() {
	slog.Info("getting tasks from the repository...")

	from := time.Now()
	quantity := s.capacity
	status := entities.StatusPending
	asc := true
	n := len(s.heap)

	if len(s.heap) > 0 {
		lastTask := s.heap[n-1]
		from = lastTask.GetRunAt()

		quantity -= n
	}

	tasks := s.taskRepository.Get(&status, &asc, &quantity, &from)

	slog.Debug(fmt.Sprintf("adding %d tasks in the scheduler...", len(tasks)))

	for _, task := range tasks {
		s.Add(&task)
	}
}

func (s *Scheduler) upsert(task *entities.Task) {
	index, ok := indexes[task.GetId()]

	slog.Debug(fmt.Sprintf("current indexes %v", indexes))

	if ok {
		s.mu.Lock()
		heap.Remove(&s.heap, index)
		heap.Push(&s.heap, task)
		s.mu.Unlock()

		slog.Info(fmt.Sprintf("task updated into index %d, task id %s", index, task.GetId()))
		return
	}

	s.mu.Lock()
	heap.Push(&s.heap, task)
	indexes[task.GetId()] = len(s.heap) - 1
	s.mu.Unlock()

	s.removeExtra()

	slog.Info(
		fmt.Sprintf(
			"new task %s added to heap and now have %d tasks on heap",
			task.GetId(),
			len(s.heap),
		),
	)
}

func (s *Scheduler) removeExtra() {
	if len(s.heap) <= s.capacity {
		return
	}

	for i := len(s.heap) - 1; i > 19; i-- {
		s.mu.Lock()
		task := heap.Remove(&s.heap, i).(*entities.Task)
		delete(indexes, task.GetId())
		s.mu.Unlock()
	}
}

func (s *Scheduler) Add(task *entities.Task) {
	s.addCh <- task
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}
