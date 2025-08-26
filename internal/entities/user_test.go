package entities_test

import (
	"scheduler/internal/errors"
	"testing"

	. "scheduler/internal/entities"

	. "scheduler/test"
)

func TestUser_New(t *testing.T) {
	user, err := NewUser("testuser", "foo@gmail.com", "Password@123")

	Ok(t, err)

	Equals(t, true, user.CheckPasswordHash("Password@123"))
	Equals(t, 0, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())
}

func TestUser_InvalidUsername(t *testing.T) {
	_, err := NewUser("test", "test@email.com", "Password@123")

	Assert(t, err != nil, "want invalid username got valid")
	Equals(t, errors.INVALID_FIELD_VALUE("username").Error(), err.Error())
}

func TestUser_InvalidEmail(t *testing.T) {
	_, err := NewUser("testuser", "test@.com", "Password@123")

	Assert(t, err != nil, "want invalid email got valid")
	Equals(t, errors.INVALID_FIELD_VALUE("email").Error(), err.Error())
}

func TestUser_InvalidPassword(t *testing.T) {
	_, err := NewUser("testuser", "test@email.com", "123456")

	Assert(t, err != nil, "want invalid password got valid")
	Equals(t, errors.INVALID_FIELD_VALUE("password").Error(), err.Error())
}

func TestUser_AddCredits(t *testing.T) {
	user, err := NewUser("testuser", "foo@gmail.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())

	user.AddCredits(10)

	Equals(t, 10, user.GetCredits())
}

func TestUser_RemoveCredits(t *testing.T) {
	user, err := NewUser("testuser", "foo@gmail.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())

	user.AddCredits(10)

	Equals(t, 10, user.GetCredits())

	err = user.RemoveCredits(10)

	Ok(t, err)
	Equals(t, 0, user.GetCredits())
}

func TestUser_AddFrozenCredits(t *testing.T) {
	user, err := NewUser("testuser", "foo@gmail.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())

	user.AddCredits(10)

	Equals(t, 10, user.GetCredits())

	err = user.AddFrozenCredits(10)

	Ok(t, err)
	Equals(t, 10, user.GetFrozenCredits())
	Equals(t, 0, user.GetCredits())
}

func TestUser_RemoveFrozenCreditsWithoutRefund(t *testing.T) {
	user, err := NewUser("testuser", "foo@gmail.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())

	user.AddCredits(10)

	Equals(t, 10, user.GetCredits())

	err = user.AddFrozenCredits(10)

	Ok(t, err)
	Equals(t, 10, user.GetFrozenCredits())
	Equals(t, 0, user.GetCredits())

	err = user.RemoveFrozenCredits(10, false)

	Ok(t, err)
	Equals(t, 0, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())
}

func TestUser_RemoveFrozenCreditsWithRefund(t *testing.T) {
	user, err := NewUser("testuser", "foo@gmail.com", "Password@123")

	Ok(t, err)

	Equals(t, 0, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())

	user.AddCredits(10)

	Equals(t, 10, user.GetCredits())

	err = user.AddFrozenCredits(10)

	Ok(t, err)
	Equals(t, 10, user.GetFrozenCredits())
	Equals(t, 0, user.GetCredits())

	err = user.RemoveFrozenCredits(10, true)

	Ok(t, err)
	Equals(t, 10, user.GetCredits())
	Equals(t, 0, user.GetFrozenCredits())
}
