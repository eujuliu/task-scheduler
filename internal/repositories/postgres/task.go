package postgres_repos

import (
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/persistence"
	"scheduler/pkg/postgres"
)

type PostgresTaskRepository struct {
	db postgres.Database
}

func NewPostgresTaskRepository() *PostgresTaskRepository {
	return &PostgresTaskRepository{
		db: *postgres.DB,
	}
}

func (r *PostgresTaskRepository) Get() []entities.Task {
	db := r.db.GetInstance()

	var tasks []persistence.TaskModel
	var result []entities.Task

	db.Find(&tasks)

	for _, task := range tasks {
		result = append(result, *persistence.ToTaskDomain(&task))
	}

	return result
}

func (r *PostgresTaskRepository) GetByUserId(userId string) []entities.Task {
	db := r.db.GetInstance()

	var tasks []persistence.TaskModel
	var result []entities.Task

	db.Find(&tasks, "user_id = ?", userId)

	for _, task := range tasks {
		result = append(result, *persistence.ToTaskDomain(&task))
	}

	return result
}

func (r *PostgresTaskRepository) GetFirstById(
	id string,
) (*entities.Task, error) {
	db := r.db.GetInstance()

	var task persistence.TaskModel

	if err := db.First(&task, "id = ?", id).Error; err != nil {
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	return persistence.ToTaskDomain(&task), nil
}

func (r *PostgresTaskRepository) GetFirstByReferenceId(
	id string,
) (*entities.Task, error) {
	db := r.db.GetInstance()

	var task persistence.TaskModel

	if err := db.First(&task, "reference_id = ?", id).Error; err != nil {
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	return persistence.ToTaskDomain(&task), nil
}

func (r *PostgresTaskRepository) GetFirstByIdempotencyKey(
	key string,
) (*entities.Task, error) {
	db := r.db.GetInstance()

	var task persistence.TaskModel

	if err := db.First(&task, "idempotency_key = ?", key).Error; err != nil {
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	return persistence.ToTaskDomain(&task), nil
}

func (r *PostgresTaskRepository) Create(task *entities.Task) error {
	db := r.db.GetInstance()

	m := persistence.ToTaskModel(task)

	err := db.Create(m).Error

	return err
}

func (r *PostgresTaskRepository) Update(task *entities.Task) error {
	db := r.db.GetInstance()

	m := persistence.ToTaskModel(task)

	err := db.Model(&m).Updates(m.ToMap()).Error

	return err
}

func (r *PostgresTaskRepository) Delete(id string) error {
	db := r.db.GetInstance()

	err := db.Delete(&persistence.TaskModel{}, "id = ?", id).Error

	return err
}
