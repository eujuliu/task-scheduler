package entities_test

import (
	"scheduler/internal/errors"
	"testing"
	"time"

	. "scheduler/internal/entities"

	"github.com/google/uuid"

	. "scheduler/test"
)

func TestPasswordRecovery_New(t *testing.T) {
	recovery, err := NewPasswordRecovery(uuid.NewString(), 5*time.Minute)

	Ok(t, err)

	Equals(t, true, recovery.InTime(time.Now()))
}

func TestPasswordRecovery_SmallerTime(t *testing.T) {
	_, err := NewPasswordRecovery(uuid.NewString(), 1*time.Minute)

	Assert(t, err != nil, "expect error got success")
	Equals(
		t,
		errors.INVALID_FIELD_VALUE("expiration time", nil).Error(),
		err.Error(),
	)
}

func TestPasswordRecovery_GreaterTime(t *testing.T) {
	_, err := NewPasswordRecovery(uuid.NewString(), 11*time.Minute)

	Assert(t, err != nil, "expect error got success")
	Equals(
		t,
		errors.INVALID_FIELD_VALUE("expiration time", nil).Error(),
		err.Error(),
	)
}
