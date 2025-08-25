package in_memory_repos

import (
	"scheduler/internal/entities"
	"slices"
)

type InMemoryTaskRepository struct {
	tasks []entities.Task
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: []entities.Task{},
	}
}

func (r *InMemoryTaskRepository) Get() []entities.Task {
	return r.tasks
}

func (r *InMemoryTaskRepository) GetByUserId(userId string) []entities.Task {
	tasks := []entities.Task{}

	for _, task := range r.tasks {
		if task.GetUserId() == userId {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

func (r *InMemoryTaskRepository) GetFirstById(
	id string,
) (*entities.Task, error) {
	for _, task := range r.tasks {
		if task.GetId() == id {
			return &task, nil
		}
	}

	return nil, nil
}

func (r *InMemoryTaskRepository) GetFirstByReferenceId(
	id string,
) (*entities.Task, error) {
	for _, task := range r.tasks {
		if task.GetReferenceId() == id {
			return &task, nil
		}
	}

	return nil, nil
}

func (r *InMemoryTaskRepository) GetFirstByIdempotencyKey(
	key string,
) (*entities.Task, error) {
	for _, task := range r.tasks {
		if task.GetIdempotencyKey() == key {
			return &task, nil
		}
	}

	return nil, nil
}

func (r *InMemoryTaskRepository) Create(task *entities.Task) error {
	r.tasks = append(r.tasks, *task)

	return nil
}

func (r *InMemoryTaskRepository) Update(task *entities.Task) error {
	for i, t := range r.tasks {
		if task.GetId() == t.GetId() {
			r.tasks[i] = *task

			return nil
		}
	}

	return nil
}

func (r *InMemoryTaskRepository) Delete(id string) error {
	index := slices.IndexFunc(r.tasks, func(t entities.Task) bool {
		return t.GetId() == id
	})

	r.tasks = slices.Delete(r.tasks, index, index+1)

	return nil
}
