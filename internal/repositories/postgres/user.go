package postgres_repos

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/persistence"
	"scheduler/pkg/postgres"
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

	return err
}

func (r *PostgresUserRepository) Update(user *entities.User) error {
	db := r.db.GetInstance()

	m, err := persistence.ToUserModel(user)
	if err != nil {
		return err
	}

	db.Model(m).Updates(m)

	return nil
}

func (r *PostgresUserRepository) Delete(id string) error {
	db := r.db.GetInstance()

	db.Delete(&persistence.UserModel{}, id)

	return nil
}
