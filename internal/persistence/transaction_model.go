package persistence

import "scheduler/internal/entities"

type TransactionModel struct {
	BaseModel
	UserID         string    `gorm:"type:uuid;not null;index"`
	User           UserModel `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Credits        int       `gorm:"not null"`
	Amount         int       `gorm:"not null"`
	Currency       string    `gorm:"type:varchar(8);not null"`
	Status         string    `gorm:"type:varchar(10);not null"`
	Type           string    `gorm:"type:varchar(15);not null"`
	ReferenceID    string    `gorm:"type:varchar(128);not null"`
	IdempotencyKey string    `gorm:"type:varchar(128);not null"`
}

func (TransactionModel) TableName() string { return "transactions" }

func ToTransactionModel(t *entities.Transaction) *TransactionModel {
	return &TransactionModel{
		BaseModel: BaseModel{
			ID:        t.GetId(),
			CreatedAt: t.GetCreatedAt(),
			UpdatedAt: t.GetUpdatedAt(),
		},
		UserID:         t.GetUserId(),
		Credits:        t.GetCredits(),
		Amount:         t.GetAmount(),
		Currency:       t.GetCurrency(),
		Status:         t.GetStatus(),
		Type:           t.GetType(),
		ReferenceID:    t.GetReferenceId(),
		IdempotencyKey: t.GetIdempotencyKey(),
	}
}

func ToTransactionDomain(m *TransactionModel) *entities.Transaction {
	return entities.HydrateTransaction(
		m.ID,
		m.UserID,
		m.Credits,
		m.Amount,
		m.Currency,
		m.Status,
		m.Type,
		m.ReferenceID,
		m.IdempotencyKey,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

func (m *TransactionModel) ToMap() map[string]any {
	return map[string]any{
		"ID":             m.ID,
		"UserID":         m.UserID,
		"Credits":        m.Credits,
		"Amount":         m.Amount,
		"Currency":       m.Currency,
		"Status":         m.Status,
		"Type":           m.Type,
		"ReferenceID":    m.ReferenceID,
		"IdempotencyKey": m.IdempotencyKey,
		"CreatedAt":      m.CreatedAt,
		"UpdatedAt":      m.UpdatedAt,
	}
}
