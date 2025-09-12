package queue

import "sync"

type InMemoryQueue struct {
	queues map[string]chan []byte
	mu     sync.RWMutex
	done   chan struct{}
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{queues: make(map[string]chan []byte), done: make(chan struct{})}
}

func (q *InMemoryQueue) ensureQueue(name string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.queues[name]; !ok {
		q.queues[name] = make(chan []byte, 100)
	}

	go func() {
		<-q.done
		close(q.queues[name])
	}()
}

func (q *InMemoryQueue) Publish(name string, data []byte) error {
	q.ensureQueue(name)
	q.queues[name] <- data
	return nil
}

func (q *InMemoryQueue) Consume(name string) (<-chan []byte, error) {
	q.ensureQueue(name)
	return q.queues[name], nil
}

func (q *InMemoryQueue) Close() {
	close(q.done)
}
