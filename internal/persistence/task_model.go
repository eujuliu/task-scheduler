package persistence

import (
	"scheduler/internal/entities"
	"time"
)

type TaskModel struct {
	BaseModel
	Type           string    `gorm:"type:varchar(50);not null"`
	UserID         string    `gorm:"type:uuid;not null;index"`
	User           UserModel `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Cost           int       `gorm:"not null"`
	Status         string    `gorm:"type:varchar(10);not null"`
	RunAt          time.Time `gorm:"not null"`
	Timezone       string    `gorm:"not null"`
	Retries        int       `gorm:"not null"`
	Priority       int       `gorm:"not null"`
	ReferenceID    string    `gorm:"type:varchar(128);not null"`
	IdempotencyKey string    `gorm:"type:varchar(128);not null"`
}

func (TaskModel) TableName() string { return "tasks" }

func ToTaskModel(t *entities.Task) *TaskModel {
	return &TaskModel{
		BaseModel: BaseModel{
			ID:        t.GetId(),
			CreatedAt: t.GetCreatedAt(),
			UpdatedAt: t.GetUpdatedAt(),
		},
		Type:           t.GetType(),
		UserID:         t.GetUserId(),
		Cost:           t.GetCost(),
		Status:         t.GetStatus(),
		RunAt:          t.GetRunAt(),
		Timezone:       t.GetTimezone(),
		Retries:        t.GetRetries(),
		Priority:       t.GetPriority(),
		ReferenceID:    t.GetReferenceId(),
		IdempotencyKey: t.GetIdempotencyKey(),
	}
}

func ToTaskDomain(m *TaskModel) *entities.Task {
	return entities.HydrateTask(
		m.ID,
		m.Type,
		m.UserID,
		m.Cost,
		m.Status,
		m.RunAt,
		m.Timezone,
		m.Retries,
		m.Priority,
		m.ReferenceID,
		m.IdempotencyKey,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

func (m *TaskModel) ToMap() map[string]any {
	return map[string]any{
		"ID":             m.ID,
		"Type":           m.Type,
		"UserID":         m.UserID,
		"Cost":           m.Cost,
		"Status":         m.Status,
		"RunAt":          m.RunAt,
		"Timezone":       m.Timezone,
		"Retries":        m.Retries,
		"Priority":       m.Priority,
		"ReferenceID":    m.ReferenceID,
		"IdempotencyKey": m.IdempotencyKey,
		"CreatedAt":      m.CreatedAt,
		"UpdatedAt":      m.UpdatedAt,
	}
}
