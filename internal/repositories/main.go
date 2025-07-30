package repos

import "scheduler/internal/entities"

type IUserRepository interface {
	Get() []entities.User
	GetById(string) (*entities.User, error)
	GetByEmail(string) (*entities.User, error)
	Create(*entities.User) bool
	Update(*entities.User) bool
	Delete(string) (bool, error)
}
