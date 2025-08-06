package entities_test

import (
	"scheduler/internal/entities"
	. "scheduler/test"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewPasswordRecovery(t *testing.T) {
	recovery, err := entities.NewPasswordRecovery(uuid.NewString(), 5*time.Minute)

	Ok(t, err)

	Equals(t, true, recovery.InTime(time.Now()))
}

func TestPasswordRecoveryWithSmallerTime(t *testing.T) {
	_, err := entities.NewPasswordRecovery(uuid.NewString(), 1*time.Minute)

	Assert(t, err != nil, "expect error got success")
}

func TestPasswordRecoveryWithGreaterTime(t *testing.T) {
	_, err := entities.NewPasswordRecovery(uuid.NewString(), 11*time.Minute)

	Assert(t, err != nil, "expect error got success")
}
