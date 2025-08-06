package entities_test

import (
	. "scheduler/internal/entities"
	. "scheduler/test"
	"testing"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("testuser", "foo@gmail.com", "Password@123")

	Ok(t, err)

	Equals(t, true, user.CheckPasswordHash("Password@123"))
}

func TestSetInvalidUsername(t *testing.T) {
	_, err := NewUser("test", "test@email.com", "Password@123")

	Assert(t, err != nil, "want invalid username got valid")
}

func TestSetInvalidEmail(t *testing.T) {
	_, err := NewUser("testuser", "test@.com", "Password@123")

	Assert(t, err != nil, "want invalid email got valid")
}

func TestSetInvalidPassword(t *testing.T) {
	_, err := NewUser("testuser", "test@email.com", "123456")

	Assert(t, err != nil, "want invalid password got valid")
}
