package entities

import (
	"scheduler/internal/errors"
	"time"
)

type PasswordRecovery struct {
	BaseEntity
	userId     string
	expiration time.Duration
}

func NewPasswordRecovery(
	userId string,
	expiration time.Duration,
) (*PasswordRecovery, error) {
	recovery := &PasswordRecovery{
		BaseEntity: *NewBaseEntity(),
		userId:     userId,
	}

	err := recovery.SetExpiration(expiration)
	if err != nil {
		return nil, err
	}

	return recovery, nil
}

func HydratePasswordRecovery(
	id, userId string,
	expiration time.Duration,
	createdAt time.Time,
	updatedAt time.Time,
) *PasswordRecovery {
	return &PasswordRecovery{
		BaseEntity: BaseEntity{
			id:        id,
			createdAt: createdAt,
			updatedAt: updatedAt,
		},

		userId:     userId,
		expiration: expiration,
	}
}

func (pr *PasswordRecovery) GetUserId() string {
	return pr.userId
}

func (pr *PasswordRecovery) SetExpiration(expiration time.Duration) error {
	if expiration < 5*time.Minute || expiration > 10*time.Minute {
		return errors.INVALID_FIELD_VALUE("expiration time")
	}

	pr.expiration = expiration

	return nil
}

func (pr *PasswordRecovery) GetExpiration() time.Duration {
	return pr.expiration
}

func (pr *PasswordRecovery) InTime(now time.Time) bool {
	expirationTime := pr.createdAt.Add(pr.expiration)

	return expirationTime.After(now)
}
