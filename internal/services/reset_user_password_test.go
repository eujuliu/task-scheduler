package services_test

import (
	. "scheduler/test"
	"testing"
)

func TestResetPassword(t *testing.T) {
	teardown := Setup(t)
	defer teardown(t)

	_, err := CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	Ok(t, err)

	recovery, err := ForgotUserPasswordService.Execute("test@email.com")

	Ok(t, err)

	err = ResetUserPasswordService.Execute(recovery.Id, "TestPassword@123")

	Ok(t, err)

	_, err = GetUserService.Execute("test@email.com", "TestPassword@123")

	Ok(t, err)
}
