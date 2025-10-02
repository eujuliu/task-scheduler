package postgres_repos

import (
	"fmt"
	"scheduler/internal/entities"
	"scheduler/internal/errors"
	"scheduler/internal/persistence"
	"scheduler/pkg/postgres"
	"time"
)

type PostgresTaskRepository struct {
	db *postgres.Database
}

func NewPostgresTaskRepository(db *postgres.Database) *PostgresTaskRepository {
	return &PostgresTaskRepository{
		db: db,
	}
}

func (r *PostgresTaskRepository) Get(
	status *string,
	asc *bool,
	limit *int,
	from *time.Time,
) []entities.Task {
	db := r.db.Get()

	var tasks []persistence.TaskModel
	var result []entities.Task

	query := db.Model(&tasks)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if from != nil {
		query = query.Where("run_at > ?", *from)
	}

	if asc != nil {
		order := "asc"

		if !*asc {
			order = "desc"
		}

		query = query.Order(fmt.Sprintf("run_at %s", order))
	}

	if limit != nil {
		query = query.Limit(*limit)
	}

	query.Find(&tasks)

	for _, task := range tasks {
		result = append(result, *persistence.ToTaskDomain(&task))
	}

	return result
}

func (r *PostgresTaskRepository) GetByUserId(
	userId string,
	offset *int,
	limit *int,
	orderBy *string,
) []entities.Task {
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

	var tasks []persistence.TaskModel
	var result []entities.Task

	db.Find(&tasks, "user_id = ?", userId).Offset(*offset).
		Limit(*limit).
		Order(fmt.Sprintf("updated_at %v", *orderBy))

	for _, task := range tasks {
		result = append(result, *persistence.ToTaskDomain(&task))
	}

	return result
}

func (r *PostgresTaskRepository) GetFirstById(
	id string,
) (*entities.Task, error) {
	db := r.db.Get()

	var task persistence.TaskModel

	if err := db.First(&task, "id = ?", id).Error; err != nil {
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	return persistence.ToTaskDomain(&task), nil
}

func (r *PostgresTaskRepository) GetFirstByReferenceId(
	id string,
) (*entities.Task, error) {
	db := r.db.Get()

	var task persistence.TaskModel

	if err := db.First(&task, "reference_id = ?", id).Error; err != nil {
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	return persistence.ToTaskDomain(&task), nil
}

func (r *PostgresTaskRepository) GetFirstByIdempotencyKey(
	key string,
) (*entities.Task, error) {
	db := r.db.Get()

	var task persistence.TaskModel

	if err := db.First(&task, "idempotency_key = ?", key).Error; err != nil {
		return nil, errors.TASK_NOT_FOUND_ERROR()
	}

	return persistence.ToTaskDomain(&task), nil
}

func (r *PostgresTaskRepository) Create(task *entities.Task) error {
	db := r.db.Get()

	m := persistence.ToTaskModel(task)

	err := db.Create(m).Error

	return err
}

func (r *PostgresTaskRepository) Update(task *entities.Task) error {
	db := r.db.Get()

	m := persistence.ToTaskModel(task)

	err := db.Model(&m).Updates(m.ToMap()).Error

	return err
}

func (r *PostgresTaskRepository) Delete(id string) error {
	db := r.db.Get()

	err := db.Delete(&persistence.TaskModel{}, "id = ?", id).Error

	return err
}
