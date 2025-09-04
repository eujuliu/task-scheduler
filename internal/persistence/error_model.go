package persistence

import (
	"scheduler/internal/entities"
)

type ErrorModel struct {
	BaseModel
	Type        string    `gorm:"type:varchar(50);not null"`
	UserID      string    `gorm:"type:uuid;not null;index"`
	User        UserModel `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Reason      string    `gorm:"type:varchar(320);not null"`
	ReferenceID string    `gorm:"type:varchar(128);not null"`
}

func (ErrorModel) TableName() string { return "errors" }

func ToErrorModel(t *entities.Error) *ErrorModel {
	return &ErrorModel{
		BaseModel: BaseModel{
			ID:        t.GetId(),
			CreatedAt: t.GetCreatedAt(),
			UpdatedAt: t.GetUpdatedAt(),
		},
		Type:        t.GetType(),
		UserID:      t.GetUserId(),
		ReferenceID: t.GetReferenceId(),
		Reason:      t.GetReason(),
	}
}

func ToErrorDomain(m *ErrorModel) *entities.Error {
	return entities.HydrateError(
		m.ID,
		m.Type,
		m.UserID,
		m.Reason,
		m.ReferenceID,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
