package services_test

import (
	. "scheduler/test"
	"testing"
	"time"
)

func TestForgotPassword(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "user@email.com", "Password@123")

	Ok(t, err)

	recovery, err := ForgotUserPasswordService.Execute(user.GetEmail())

	Ok(t, err)

	Equals(t, true, recovery.InTime(time.Now()))
}

func TestExpiredToken(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	user, err := CreateUserService.Execute("testuser", "user@email.com", "Password@123")

	Ok(t, err)

	recovery, err := ForgotUserPasswordService.Execute(user.GetEmail())

	Ok(t, err)

	Equals(t, false, recovery.InTime(time.Now().Add(6*time.Minute)))
}
