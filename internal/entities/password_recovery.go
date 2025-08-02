package entities

import (
	"scheduler/internal/errors"
	"time"

	"github.com/google/uuid"
)

type PasswordRecovery struct {
	Id             string
	UserId         string
	CreatedAt      time.Time
	ExpirationTime time.Duration
}

func NewPasswordRecovery(userId string, expirationTime time.Duration) (*PasswordRecovery, error) {
	if !validExpirationTime(expirationTime) {
		return nil, errors.EXPIRATION_TIME_INVALID()
	}

	return &PasswordRecovery{
		Id:             uuid.NewString(),
		UserId:         userId,
		CreatedAt:      time.Now(),
		ExpirationTime: expirationTime,
	}, nil
}

func validExpirationTime(expiration time.Duration) bool {
	return expiration >= 5*time.Minute && expiration <= 10*time.Minute
}

func (r *PasswordRecovery) InTime(now time.Time) bool {
	expirationTime := r.CreatedAt.Add(r.ExpirationTime)

	return expirationTime.After(now)
}
