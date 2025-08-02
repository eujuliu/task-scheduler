package repos

import "scheduler/internal/entities"

type IUserRepository interface {
	Get() []entities.User
	GetById(string) (*entities.User, error)
	GetByEmail(string) (*entities.User, error)
	Create(*entities.User) error
	Update(*entities.User) error
	Delete(string) error
}

type IPasswordRecoveryRepository interface {
	Get() []entities.PasswordRecovery
	GetById(string) (*entities.PasswordRecovery, error)
	GetByUserId(string) (*entities.PasswordRecovery, error)
	Create(*entities.PasswordRecovery) error
	Delete(string) error
}
