package persistence

import (
	"scheduler/internal/entities"
)

type UserModel struct {
	BaseModel
	Username      string `gorm:"size:20;uniqueIndex;not null"`
	Email         string `gorm:"size:255;uniqueIndex;not null"`
	Password      string `gorm:"size:72;not null"`
	Credits       int    `gorm:"not null"`
	FrozenCredits int    `gorm:"not null"`
}

func (UserModel) TableName() string { return "users" }

func ToUserModel(u *entities.User) (*UserModel, error) {
	pw, err := u.GetPassword()
	if err != nil {
		return nil, err
	}

	return &UserModel{
		BaseModel: BaseModel{
			ID:        u.GetId(),
			CreatedAt: u.GetCreatedAt(),
			UpdatedAt: u.GetUpdatedAt(),
		},

		Username:      u.GetUsername(),
		Email:         u.GetEmail(),
		Password:      pw,
		Credits:       u.GetCredits(),
		FrozenCredits: u.GetFrozenCredits(),
	}, nil
}

func ToUserDomain(m *UserModel) *entities.User {
	return entities.HydrateUser(
		m.ID,
		m.CreatedAt,
		m.UpdatedAt,
		m.Username,
		m.Email,
		m.Password,
		m.Credits,
		m.FrozenCredits,
	)
}

func (m *UserModel) ToMap() map[string]any {
	return map[string]any{
		"ID":            m.ID,
		"CreatedAt":     m.CreatedAt,
		"UpdatedAt":     m.UpdatedAt,
		"Username":      m.Username,
		"Email":         m.Email,
		"Password":      m.Password,
		"Credits":       m.Credits,
		"FrozenCredits": m.FrozenCredits,
	}
}
