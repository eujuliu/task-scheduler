package in_memory_repos

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"slices"
)

type InMemoryUserRepository struct {
	users []entities.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: []entities.User{},
	}
}

func (r *InMemoryUserRepository) Get() []entities.User {
	return r.users
}

func (r *InMemoryUserRepository) GetById(id string) (*entities.User, error) {
	for _, v := range r.users {
		if v.GetId() == id {
			return &v, nil
		}
	}

	return nil, errors.USER_NOT_FOUND_ERROR()
}

func (r *InMemoryUserRepository) GetByEmail(email string) (*entities.User, error) {
	for _, v := range r.users {
		if v.GetEmail() == email {
			return &v, nil
		}
	}

	return nil, errors.USER_NOT_FOUND_ERROR()
}

func (r *InMemoryUserRepository) Create(user *entities.User) error {
	r.users = append(r.users, *user)

	return nil
}

func (r *InMemoryUserRepository) Update(user *entities.User) error {
	for i, u := range r.users {
		if user.GetId() == u.GetId() {
			r.users[i] = *user

			return nil
		}
	}

	return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
	index := slices.IndexFunc(r.users, func(u entities.User) bool {
		return u.GetId() == id
	})

	r.users = slices.Delete(r.users, index, index+1)

	return nil
}
