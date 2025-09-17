package scheduler

import "scheduler/internal/entities"

type Heap []*entities.Task

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	if h[i].GetRunAt().Equal(h[j].GetRunAt()) {
		return h[i].GetPriority() > h[j].GetPriority()
	}

	return h[i].GetRunAt().Before(h[j].GetRunAt())
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *Heap) Push(x any) {
	*h = append(*h, x.(*entities.Task))
}

func (h *Heap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *Heap) Peek() any {
	old := *h
	n := len(old)
	x := old[n-1]
	return x
}
