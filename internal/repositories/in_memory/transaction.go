package in_memory_repos

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"slices"
)

type InMemoryTransactionRepository struct {
	transactions []entities.Transaction
}

func NewInMemoryTransactionRepository() *InMemoryTransactionRepository {
	return &InMemoryTransactionRepository{
		transactions: []entities.Transaction{},
	}
}

func (repo *InMemoryTransactionRepository) Get() []entities.Transaction {
	return repo.transactions
}

func (repo *InMemoryTransactionRepository) GetByUserId(userId string) []entities.Transaction {
	var result = []entities.Transaction{}

	for _, transaction := range repo.transactions {
		if transaction.GetUserId() == userId {
			result = append(result, transaction)
		}
	}

	return result
}

func (repo *InMemoryTransactionRepository) GetFirstById(id string) (*entities.Transaction, error) {
	for _, transaction := range repo.transactions {
		if transaction.GetId() == id {
			return &transaction, nil
		}
	}

	return nil, errors.TRANSACTION_NOT_FOUND()
}

func (repo *InMemoryTransactionRepository) GetFirstByReferenceId(id string) (*entities.Transaction, error) {
	for _, transaction := range repo.transactions {
		if transaction.GetReferenceId() == id {
			return &transaction, nil
		}
	}

	return nil, errors.TRANSACTION_NOT_FOUND()
}

func (repo *InMemoryTransactionRepository) GetFirstByIdempotencyKey(key string) (*entities.Transaction, error) {
	for _, transaction := range repo.transactions {
		if transaction.GetIdempotencyKey() == key {
			return &transaction, nil
		}
	}

	return nil, errors.TRANSACTION_NOT_FOUND()
}

func (repo *InMemoryTransactionRepository) Create(transaction *entities.Transaction) error {
	repo.transactions = append(repo.transactions, *transaction)

	return nil
}

func (repo *InMemoryTransactionRepository) Update(transaction *entities.Transaction) error {
	for i, t := range repo.transactions {
		if transaction.GetId() == t.GetId() {
			repo.transactions[i] = *transaction

			return nil
		}
	}

	return nil
}

func (repo *InMemoryTransactionRepository) Delete(id string) error {
	index := slices.IndexFunc(repo.transactions, func(t entities.Transaction) bool {
		return t.GetId() == id
	})

	repo.transactions = slices.Delete(repo.transactions, index, index+1)

	return nil
}
