package in_memory_repos

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"slices"
)

type InMemoryUserRepository struct {
	users []entities.User
}

func NewUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: []entities.User{},
	}
}

func (r *InMemoryUserRepository) Get() []entities.User {
	return r.users
}

func (r *InMemoryUserRepository) GetById(id string) (*entities.User, error) {
	for _, v := range r.users {
		if v.Id == id {
			return &v, nil
		}
	}

	return nil, errors.USER_NOT_FOUND_ERROR()
}

func (r *InMemoryUserRepository) GetByEmail(email string) (*entities.User, error) {
	for _, v := range r.users {
		if v.Email == email {
			return &v, nil
		}
	}

	return nil, errors.USER_NOT_FOUND_ERROR()
}

func (r *InMemoryUserRepository) Create(user *entities.User) (bool, error) {
	r.users = append(r.users, *user)

	return true, nil
}

func (r *InMemoryUserRepository) Update(user *entities.User) (bool, error) {
	for _, u := range r.users {
		if user.Id == u.Id {
			u = *user

			return true, nil
		}
	}

	return false, nil
}

func (r *InMemoryUserRepository) Delete(id string) (bool, error) {
	index := slices.IndexFunc(r.users, func(u entities.User) bool {
		return u.Id == id
	})

	r.users = slices.Delete(r.users, index, index+1)

	return true, nil
}
