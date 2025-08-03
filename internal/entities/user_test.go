package entities_test

import (
	. "scheduler/internal/entities"
	. "scheduler/test"
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("testuser", "foo@gmail.com", "Password@123")

	Ok(t, err)

	Ok(t, uuid.Validate(user.Id))

	Equals(t, true, user.CheckPasswordHash("Password@123"))
}

func TestSetPassword(t *testing.T) {
	user, err := NewUser("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	err = user.SetPassword("TestPassword@123")

	Ok(t, err)

	Equals(t, true, user.CheckPasswordHash("TestPassword@123"))
}

func TestUserInvalidUsername(t *testing.T) {
	_, err := NewUser("test", "test@email.com", "Password@123")

	Assert(t, err != nil, "want invalid username got valid")
}

func TestUserInvalidEmail(t *testing.T) {
	_, err := NewUser("testuser", "test@.com", "Password@123")

	Assert(t, err != nil, "want invalid email got valid")
}

func TestUserInvalidPassword(t *testing.T) {
	_, err := NewUser("testuser", "test@email.com", "123456")

	Assert(t, err != nil, "want invalid password got valid")
}
