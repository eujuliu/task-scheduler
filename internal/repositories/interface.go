package repos

import "scheduler/internal/entities"

type IUserRepository interface {
	Get() []entities.User
	GetFirstById(string) (*entities.User, error)
	GetFirstByEmail(string) (*entities.User, error)
	Create(*entities.User) error
	Update(*entities.User) error
	Delete(string) error
}

type IPasswordRecoveryRepository interface {
	Get() []entities.PasswordRecovery
	GetFirstById(string) (*entities.PasswordRecovery, error)
	GetFirstByUserId(string) (*entities.PasswordRecovery, error)
	Create(*entities.PasswordRecovery) error
	Delete(id string) error
}

type ITransactionRepository interface {
	Get() []entities.Transaction
	GetByUserId(userId string) []entities.Transaction
	GetFirstById(id string) (*entities.Transaction, error)
	GetFirstByReferenceId(id string) (*entities.Transaction, error)
	GetFirstByIdempotencyKey(key string) (*entities.Transaction, error)
	Create(*entities.Transaction) error
	Update(*entities.Transaction) error
	Delete(id string) error
}

type IErrorRepository interface {
	Get() []entities.Error
	GetByUserId(userId string) []entities.Error
	GetFirstByUserId(userId string) (*entities.Error, error)
	GetFirstByReferenceId(referenceId string) (*entities.Error, error)
	Create(*entities.Error) error
}
