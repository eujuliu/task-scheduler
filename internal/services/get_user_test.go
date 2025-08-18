package services_test

import (
	"scheduler/internal/errors"
	. "scheduler/test"
	"testing"
)

func TestGetUserService(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	_, err = GetUserService.Execute("test@email.com", "Password@123")

	Ok(t, err)
}

func TestGetUserService_WrongPassword(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	_, err = GetUserService.Execute("test@email.com", "Password@12")

	Equals(t, errors.WRONG_LOGIN_DATA_ERROR().Error(), err.Error())
}
