package postgres_repos

import (
	"fmt"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/persistence"
	"scheduler/pkg/postgres"
)

type PostgresTransactionRepository struct {
	db *postgres.Database
}

func NewPostgresTransactionRepository(db *postgres.Database) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{
		db: db,
	}
}

func (r *PostgresTransactionRepository) Get() []entities.Transaction {
	db := r.db.Get()

	var transactions []persistence.TransactionModel
	var result []entities.Transaction

	db.Find(&transactions)

	for _, transaction := range transactions {
		result = append(result, *persistence.ToTransactionDomain(&transaction))
	}

	return result
}

func (r *PostgresTransactionRepository) GetByUserId(
	userId string, offset *int, limit *int, orderBy *string,
) []entities.Transaction {
	if offset == nil {
		*offset = 0
	}

	if limit == nil {
		*limit = 10
	}

	if orderBy == nil {
		*orderBy = "ASC"
	}

	db := r.db.Get()

	var transactions []persistence.TransactionModel
	var result []entities.Transaction

	db.Find(&transactions, "user_id = ?", userId).
		Offset(*offset).
		Limit(*limit).
		Order(fmt.Sprintf("updated_at %v", *orderBy))

	for _, transaction := range transactions {
		result = append(result, *persistence.ToTransactionDomain(&transaction))
	}

	return result
}

func (r *PostgresTransactionRepository) GetFirstById(
	id string,
) (*entities.Transaction, error) {
	db := r.db.Get()

	var transaction persistence.TransactionModel

	if err := db.First(&transaction, "id = ?", id).Error; err != nil {
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	return persistence.ToTransactionDomain(&transaction), nil
}

func (r *PostgresTransactionRepository) GetFirstByReferenceId(
	id string,
) (*entities.Transaction, error) {
	db := r.db.Get()

	var transaction persistence.TransactionModel

	if err := db.First(&transaction, "reference_id = ?", id).Error; err != nil {
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	return persistence.ToTransactionDomain(&transaction), nil
}

func (r *PostgresTransactionRepository) GetFirstByIdempotencyKey(
	key string,
) (*entities.Transaction, error) {
	db := r.db.Get()

	var transaction persistence.TransactionModel

	if err := db.First(&transaction, "idempotency_key = ?", key).Error; err != nil {
		return nil, errors.TRANSACTION_NOT_FOUND()
	}

	return persistence.ToTransactionDomain(&transaction), nil
}

func (r *PostgresTransactionRepository) Create(
	transaction *entities.Transaction,
) error {
	db := r.db.Get()

	m := persistence.ToTransactionModel(transaction)

	err := db.Create(m).Error

	return err
}

func (r *PostgresTransactionRepository) Update(
	transaction *entities.Transaction,
) error {
	db := r.db.Get()

	m := persistence.ToTransactionModel(transaction)

	err := db.Model(&m).Updates(m.ToMap()).Error

	return err
}
