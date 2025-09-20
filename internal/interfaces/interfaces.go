package interfaces

import (
	"context"
	"scheduler/internal/entities"
	"time"
)

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
}

type IErrorRepository interface {
	Get() []entities.Error
	GetByUserId(userId string) []entities.Error
	GetFirstByUserId(userId string) (*entities.Error, error)
	GetFirstByReferenceId(referenceId string) (*entities.Error, error)
	Create(*entities.Error) error
}

type ITaskRepository interface {
	Get(status *string, asc *bool, limit *int, from *time.Time) []entities.Task
	GetByUserId(userId string) []entities.Task
	GetFirstById(id string) (*entities.Task, error)
	GetFirstByReferenceId(id string) (*entities.Task, error)
	GetFirstByIdempotencyKey(key string) (*entities.Task, error)
	Create(*entities.Task) error
	Update(*entities.Task) error
	Delete(id string) error
}

type ICustomerPaymentGateway interface {
	GetFirstByEmail(email string) (*string, error)
	Create(internalID, username, email string, props *map[string]string) (string, error)
}

type IPaymentPaymentGateway interface {
	Create(
		internalID string,
		customerID *string,
		amount int,
		currency, idempotencyKey string,
		props *map[string]string,
	) (string, error)
}

type IQueue interface {
	Publish(key string, exchangeName string, data []byte) error
	Consume(ctx context.Context, queue string, handler func(map[string]any) error) error
}
