package postgres_repos

import (
	errorr "errors"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/persistence"
	"scheduler/pkg/postgres"

	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *postgres.Database
}

func NewPostgresUserRepository() *PostgresUserRepository {
	return &PostgresUserRepository{
		db: postgres.DB,
	}
}

func (r *PostgresUserRepository) Get() []entities.User {
	db := r.db.GetInstance()

	var users []persistence.UserModel
	var result []entities.User

	db.Find(&users)

	for _, user := range users {
		result = append(result, *persistence.ToUserDomain(&user))
	}

	return result
}

func (r *PostgresUserRepository) GetFirstById(
	id string,
) (*entities.User, error) {
	db := r.db.GetInstance()

	var user persistence.UserModel

	if err := db.First(&user, "id = ?", id).Error; err != nil {
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	return persistence.ToUserDomain(&user), nil
}

func (r *PostgresUserRepository) GetFirstByEmail(
	email string,
) (*entities.User, error) {
	db := r.db.GetInstance()

	var user persistence.UserModel

	if err := db.First(&user, "email = ?", email).Error; err != nil {
		return nil, errors.USER_NOT_FOUND_ERROR()
	}

	return persistence.ToUserDomain(&user), nil
}

func (r *PostgresUserRepository) Create(user *entities.User) error {
	db := r.db.GetInstance()

	m, err := persistence.ToUserModel(user)
	if err != nil {
		return err
	}

	err = db.Create(m).Error

	if errorr.Is(err, gorm.ErrDuplicatedKey) {
		return errors.USER_ALREADY_EXISTS_ERROR()
	}

	return err
}

func (r *PostgresUserRepository) Update(user *entities.User) error {
	db := r.db.GetInstance()

	m, err := persistence.ToUserModel(user)
	if err != nil {
		return err
	}

	err = db.Model(&m).Updates(m.ToMap()).Error

	return err
}

func (r *PostgresUserRepository) Delete(id string) error {
	db := r.db.GetInstance()

	err := db.Delete(&persistence.UserModel{}, "id = ?", id).Error

	return err
}
