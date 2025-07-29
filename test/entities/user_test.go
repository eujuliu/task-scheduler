package main

import (
	"scheduler/internal/entities"
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	user, err := entities.NewUser("julio1234", "foo@gmail.com", "Password@123")

	if err != nil {
		t.Errorf("got %v want no error!", err)
	}

	if err := uuid.Validate(user.Id); err != nil {
		t.Errorf("got %v want valid uuid!", err)
	}

	if !user.CheckPasswordHash("Password@123") {
		t.Errorf("got invalid hashed password")
	}
}

func TestUserInvalidFields(t *testing.T) {
	_, err := entities.NewUser("test", "foo@gmail.com", "Password@123")

	if err == nil {
		t.Errorf("got user want invalid username")
	}

	_, err = entities.NewUser("julio1234", "foo@.com", "Password@123")

	if err == nil {
		t.Errorf("got user want invalid email")
	}

	_, err = entities.NewUser("julio1234", "foo@gmail.com", "123456")

	if err == nil {
		t.Errorf("got user want invalid password")
	}
}
