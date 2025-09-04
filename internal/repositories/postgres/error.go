package postgres_repos

import (
	"scheduler/internal/entities"
	"scheduler/internal/persistence"
	"scheduler/pkg/postgres"
)

type PostgresErrorRepository struct {
	db postgres.Database
}

func NewPostgresErrorRepository() *PostgresErrorRepository {
	return &PostgresErrorRepository{
		db: *postgres.DB,
	}
}

func (r *PostgresErrorRepository) Get() []entities.Error {
	db := r.db.GetInstance()

	var troubles []persistence.ErrorModel
	var result []entities.Error

	db.Find(&troubles)

	for _, trouble := range troubles {
		result = append(result, *persistence.ToErrorDomain(&trouble))
	}

	return result
}

func (r *PostgresErrorRepository) GetByUserId(userId string) []entities.Error {
	db := r.db.GetInstance()

	var troubles []persistence.ErrorModel
	var result []entities.Error

	db.Find(&troubles, "user_id = ?", userId)

	for _, trouble := range troubles {
		result = append(result, *persistence.ToErrorDomain(&trouble))
	}

	return result
}

func (r *PostgresErrorRepository) GetFirstByUserId(
	userId string,
) (*entities.Error, error) {
	db := r.db.GetInstance()

	var trouble persistence.ErrorModel

	if err := db.First(&trouble, "user_id = ?", userId).Error; err != nil {
		return nil, err
	}

	return persistence.ToErrorDomain(&trouble), nil
}

func (r *PostgresErrorRepository) GetFirstByReferenceId(
	id string,
) (*entities.Error, error) {
	db := r.db.GetInstance()

	var trouble persistence.ErrorModel

	if err := db.First(&trouble, "reference_id = ?", id).Error; err != nil {
		return nil, err
	}

	return persistence.ToErrorDomain(&trouble), nil
}

func (r *PostgresErrorRepository) Create(trouble *entities.Error) error {
	db := r.db.GetInstance()

	m := persistence.ToErrorModel(trouble)

	err := db.Create(m).Error

	return err
}
