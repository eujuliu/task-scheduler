package entities_test

import (
	. "scheduler/internal/entities"
	"scheduler/internal/errors"
	. "scheduler/test"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordRecovery_New(t *testing.T) {
	recovery, err := NewPasswordRecovery(uuid.NewString(), 5*time.Minute)

	Ok(t, err)

	Equals(t, true, recovery.InTime(time.Now()))
}

func TestPasswordRecovery_SmallerTime(t *testing.T) {
	_, err := NewPasswordRecovery(uuid.NewString(), 1*time.Minute)

	Assert(t, err != nil, "expect error got success")
	Equals(t, errors.INVALID_FIELD_VALUE("expiration time").Error(), err.Error())
}

func TestPasswordRecovery_GreaterTime(t *testing.T) {
	_, err := NewPasswordRecovery(uuid.NewString(), 11*time.Minute)

	Assert(t, err != nil, "expect error got success")
	Equals(t, errors.INVALID_FIELD_VALUE("expiration time").Error(), err.Error())
}
