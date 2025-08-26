package in_memory_repos

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
)

type InMemoryTransactionRepository struct {
	transactions []entities.Transaction
}

func NewInMemoryTransactionRepository() *InMemoryTransactionRepository {
	return &InMemoryTransactionRepository{
		transactions: []entities.Transaction{},
	}
}

func (r *InMemoryTransactionRepository) Get() []entities.Transaction {
	return r.transactions
}

func (r *InMemoryTransactionRepository) GetByUserId(
	userId string,
) []entities.Transaction {
	result := []entities.Transaction{}

	for _, transaction := range r.transactions {
		if transaction.GetUserId() == userId {
			result = append(result, transaction)
		}
	}

	return result
}

func (r *InMemoryTransactionRepository) GetFirstById(
	id string,
) (*entities.Transaction, error) {
	for _, transaction := range r.transactions {
		if transaction.GetId() == id {
			return &transaction, nil
		}
	}

	return nil, errors.TRANSACTION_NOT_FOUND()
}

func (r *InMemoryTransactionRepository) GetFirstByReferenceId(
	id string,
) (*entities.Transaction, error) {
	for _, transaction := range r.transactions {
		if transaction.GetReferenceId() == id {
			return &transaction, nil
		}
	}

	return nil, errors.TRANSACTION_NOT_FOUND()
}

func (r *InMemoryTransactionRepository) GetFirstByIdempotencyKey(
	key string,
) (*entities.Transaction, error) {
	for _, transaction := range r.transactions {
		if transaction.GetIdempotencyKey() == key {
			return &transaction, nil
		}
	}

	return nil, errors.TRANSACTION_NOT_FOUND()
}

func (r *InMemoryTransactionRepository) Create(
	transaction *entities.Transaction,
) error {
	r.transactions = append(r.transactions, *transaction)

	return nil
}

func (r *InMemoryTransactionRepository) Update(
	transaction *entities.Transaction,
) error {
	for i, t := range r.transactions {
		if transaction.GetId() == t.GetId() {
			r.transactions[i] = *transaction

			return nil
		}
	}

	return nil
}
