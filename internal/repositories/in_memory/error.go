package in_memory_repos

import (
	"fmt"
	"scheduler/internal/entities"
)

type InMemoryErrorRepository struct {
	errors []entities.Error
}

func NewInMemoryErrorRepository() *InMemoryErrorRepository {
	return &InMemoryErrorRepository{
		errors: []entities.Error{},
	}
}

func (r *InMemoryErrorRepository) Get() []entities.Error {
	return r.errors
}

func (r *InMemoryErrorRepository) GetByUserId(userId string) []entities.Error {
	var result []entities.Error

	for _, error := range r.errors {
		if error.GetUserId() == userId {
			result = append(result, error)
		}
	}

	return result
}

func (r *InMemoryErrorRepository) GetFirstByUserId(
	userId string,
) (*entities.Error, error) {
	for _, error := range r.errors {
		if error.GetUserId() == userId {
			return &error, nil
		}
	}

	return nil, nil
}

func (r *InMemoryErrorRepository) GetFirstByReferenceId(
	referenceId string,
) (*entities.Error, error) {
	for _, error := range r.errors {
		if error.GetReferenceId() == referenceId {
			return &error, nil
		}
	}

	return nil, nil
}

func (r *InMemoryErrorRepository) Create(err *entities.Error) error {
	fmt.Print("Test append error")
	r.errors = append(r.errors, *err)

	return nil
}

func (r *InMemoryErrorRepository) Clear() {
	r.errors = []entities.Error{}
}
