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

	"github.com/jonboulle/clockwork"
)

type Scheduler struct {
	heap           Heap
	capacity       int
	mu             sync.Mutex
	clock          clockwork.Clock
	addCh          chan *entities.Task
	updateCh       chan *entities.Task
	stopCh         chan struct{}
	values         map[string]int
	queue          interfaces.IQueue
	taskRepository interfaces.ITaskRepository
}

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
		updateCh:       make(chan *entities.Task),
		stopCh:         make(chan struct{}),
		values:         make(map[string]int),
		queue:          queue,
		taskRepository: taskRepository,
	}

	heap.Init(&sc.heap)

	sc.loadNextTasks()
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
				s.push(x)
			case <-s.stopCh:
				return
			}
		} else {
			next := s.heap[0]
			now := s.clock.Now()
			waitDuration := next.GetRunAt().Sub(now)

			if waitDuration <= 0 {
				s.mu.Lock()
				x := heap.Pop(&s.heap).(*entities.Task)
				s.mu.Unlock()

				s.ExecuteTask(x)
			} else {
				timer := s.clock.After(waitDuration)

				select {
				case <-timer:
					s.mu.Lock()
					task := heap.Pop(&s.heap).(*entities.Task)
					s.mu.Unlock()
					s.ExecuteTask(task)
				case task := <-s.addCh:
					s.push(task)
				case task := <-s.updateCh:
					s.update(task)
				case <-s.stopCh:
					return
				}
			}
		}
	}
}

func (s *Scheduler) Add(task *entities.Task) {
	s.addCh <- task
}

func (s *Scheduler) Update(task *entities.Task) {
	s.updateCh <- task
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
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
	delete(s.values, task.GetId())

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

	s.loadNextTasks()
}

func (s *Scheduler) push(task *entities.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.values[task.GetId()]; ok {
		return
	}

	if len(s.heap) < s.capacity {
		heap.Push(&s.heap, task)
		s.values[task.GetId()] = 1
		slog.Debug(fmt.Sprintf("task %s added into the heap", task.GetId()))
		return
	}

	n := len(s.heap)
	lastTask := s.heap[n-1]

	if task.GetRunAt().After(lastTask.GetRunAt()) {
		return
	}

	if task.GetRunAt().Equal(lastTask.GetRunAt()) && task.GetPriority() < lastTask.GetPriority() {
		return
	}

	heap.Remove(&s.heap, n-1)
	heap.Push(&s.heap, task)
	s.values[task.GetId()] = 1
	slog.Debug(fmt.Sprintf("task %s added into the heap", task.GetId()))
}

func (s *Scheduler) update(task *entities.Task) {
	s.mu.Lock()
	indexes := make(map[string]int)

	for i, t := range s.heap {
		indexes[t.GetId()] = i
	}
	s.mu.Unlock()

	if _, ok := indexes[task.GetId()]; !ok {
		s.push(task)

		return
	}

	s.mu.Lock()
	heap.Remove(&s.heap, indexes[task.GetId()])
	heap.Push(&s.heap, task)
	s.mu.Unlock()

	slog.Debug(fmt.Sprintf("task %s updated into the heap", task.GetId()))
}

func (s *Scheduler) loadNextTasks() {
	slog.Info("getting tasks from the repository...")

	from := s.clock.Now()
	n := len(s.heap)
	quantity := min(20, s.capacity-n)
	status := entities.StatusPending
	asc := true

	tasks := s.taskRepository.Get(&status, &asc, &quantity, &from)

	slog.Info(fmt.Sprintf("find new %d tasks for add in the heap...", len(tasks)))

	go func() {
		for _, task := range tasks {
			s.Add(&task)
		}
	}()
}
