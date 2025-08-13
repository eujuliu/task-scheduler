package services_test

import (
	"scheduler/internal/errors"
	. "scheduler/test"
	"testing"
)

func TestCreateUserService(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)
}

func TestCreateUserService_Duplicity(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	_, err = CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Equals(t, errors.USER_ALREADY_EXISTS_ERROR().Error(), err.Error())
}
