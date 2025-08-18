package entities_test

import (
	. "scheduler/internal/entities"
	. "scheduler/test"
	"testing"
)

func TestUser_New(t *testing.T) {
	user, err := NewUser("testuser", "foo@gmail.com", "Password@123")

	Ok(t, err)

	Equals(t, true, user.CheckPasswordHash("Password@123"))
}

func TestUser_InvalidUsername(t *testing.T) {
	_, err := NewUser("test", "test@email.com", "Password@123")

	Assert(t, err != nil, "want invalid username got valid")
}

func TestUser_InvalidEmail(t *testing.T) {
	_, err := NewUser("testuser", "test@.com", "Password@123")

	Assert(t, err != nil, "want invalid email got valid")
}

func TestUser_InvalidPassword(t *testing.T) {
	_, err := NewUser("testuser", "test@email.com", "123456")

	Assert(t, err != nil, "want invalid password got valid")
}
