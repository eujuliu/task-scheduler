package services_test

import (
	"testing"
	"time"

	. "scheduler/test"
)

func TestForgotPasswordService(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute(
		"testuser",
		"user@email.com",
		"Password@123",
	)

	Ok(t, err)

	recovery, err := ForgotUserPasswordService.Execute(user.GetEmail())

	Ok(t, err)

	Equals(t, true, recovery.InTime(time.Now()))
}

func TestForgotPasswordService_ExpiredToken(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute(
		"testuser",
		"user@email.com",
		"Password@123",
	)

	Ok(t, err)

	recovery, err := ForgotUserPasswordService.Execute(user.GetEmail())

	Ok(t, err)

	Equals(t, false, recovery.InTime(time.Now().Add(6*time.Minute)))
}
