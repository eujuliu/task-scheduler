package postgres_repos

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/persistence"
	"scheduler/pkg/postgres"
)

type PostgresPasswordRecoveryRepository struct {
	db *postgres.Database
}

func NewPostgresPasswordRepository() *PostgresPasswordRecoveryRepository {
	return &PostgresPasswordRecoveryRepository{
		db: postgres.DB,
	}
}

func (r *PostgresPasswordRecoveryRepository) Get() []entities.PasswordRecovery {
	db := r.db.GetInstance()

	var tokens []persistence.PasswordRecoveryModel
	var result []entities.PasswordRecovery

	db.Find(&tokens)

	for _, pr := range tokens {
		result = append(result, *persistence.ToPasswordRecoveryDomain(&pr))
	}

	return result
}

func (r *PostgresPasswordRecoveryRepository) GetFirstById(
	id string,
) (*entities.PasswordRecovery, error) {
	var token persistence.PasswordRecoveryModel

	db := r.db.GetInstance()

	if err := db.First(&token, "id = ?", id).Error; err != nil {
		return nil, errors.RECOVERY_TOKEN_NOT_FOUND()
	}

	return persistence.ToPasswordRecoveryDomain(&token), nil
}

func (r *PostgresPasswordRecoveryRepository) GetFirstByUserId(
	userId string,
) (*entities.PasswordRecovery, error) {
	var token persistence.PasswordRecoveryModel

	db := r.db.GetInstance()

	if err := db.First(&token, "userID = ?", userId).Error; err != nil {
		return nil, errors.RECOVERY_TOKEN_NOT_FOUND()
	}

	return persistence.ToPasswordRecoveryDomain(&token), nil
}

func (r *PostgresPasswordRecoveryRepository) Create(
	token *entities.PasswordRecovery,
) error {
	db := r.db.GetInstance()

	m, err := persistence.ToPasswordRecoveryModel(token)
	if err != nil {
		return err
	}

	err = db.Create(m).Error

	return err
}

func (r *PostgresPasswordRecoveryRepository) Delete(id string) error {
	db := r.db.GetInstance()

	err := db.Delete(&persistence.PasswordRecoveryModel{}, "id = ? ", id).Error

	return err
}
