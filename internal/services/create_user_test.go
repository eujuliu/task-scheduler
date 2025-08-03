package services_test

import (
	. "scheduler/test"
	"testing"
)

func TestCreateUser(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)
}

func TestDuplicity(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	_, err = CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Assert(t, err != nil, "want error and got success")
}
