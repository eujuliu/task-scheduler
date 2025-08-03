package services_test

import (
	. "scheduler/test"
	"testing"
)

func TestGetUser(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	_, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)
}

func TestGetUserWithWrongPassword(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	_, err = GetUserService.Execute("test@email.com", "Password@12")

	Assert(t, err != nil, "expect error got user")
}
