package persistence

import (
	"scheduler/internal/entities"
	"time"
)

type PasswordRecoveryModel struct {
	BaseModel
	UserID     string        `gorm:"type:uuid;not null;index"`
	User       UserModel     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Expiration time.Duration `gorm:"not null"`
}

func (PasswordRecoveryModel) TableName() string {
	return "password_recoveries"
}

func ToPasswordRecoveryModel(pr *entities.PasswordRecovery) (*PasswordRecoveryModel, error) {
	return &PasswordRecoveryModel{
		BaseModel: BaseModel{
			ID:        pr.GetId(),
			CreatedAt: pr.GetCreatedAt(),
			UpdatedAt: pr.GetUpdatedAt(),
		},

		UserID:     pr.GetUserId(),
		Expiration: pr.GetExpiration(),
	}, nil
}

func ToPasswordRecoveryDomain(m *PasswordRecoveryModel) *entities.PasswordRecovery {
	return entities.HydratePasswordRecovery(m.ID, m.UserID, m.Expiration, m.CreatedAt, m.UpdatedAt)
}
