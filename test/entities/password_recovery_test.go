package main

import (
	"scheduler/internal/entities"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewPasswordRecovery(t *testing.T) {
	recovery, err := entities.NewPasswordRecovery(uuid.NewString(), 5*time.Minute)

	if err != nil {
		t.Errorf("got error %v want recovery token", err)
	}

	if err := uuid.Validate(recovery.Id); err != nil {
		t.Errorf("got error %v want valid uuid", err)
	}

	if !recovery.InTime(time.Now()) {
		t.Error("got not in time want in time")
	}
}

func TestPasswordRecoveryWithGreaterTime(t *testing.T) {
	_, err := entities.NewPasswordRecovery(uuid.NewString(), 11*time.Minute)

	if err == nil {
		t.Error("got success want error")
	}
}
