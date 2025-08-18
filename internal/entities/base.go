package entities

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	id        string
	createdAt time.Time
	updatedAt time.Time
}

func NewBaseEntity() *BaseEntity {
	now := time.Now()

	return &BaseEntity{
		id:        uuid.NewString(),
		createdAt: now,
		updatedAt: now,
	}
}

func (b *BaseEntity) GetId() string {
	return b.id
}

func (b *BaseEntity) GetCreatedAt() time.Time {
	return b.createdAt
}

func (b *BaseEntity) SetUpdatedAt() {
	b.updatedAt = time.Now()
}

func (b *BaseEntity) GetUpdatedAt() time.Time {
	return b.updatedAt
}
