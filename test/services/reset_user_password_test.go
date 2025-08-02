package main

import (
	"scheduler/test"
	"testing"
)

func TestResetPassword(t *testing.T) {
	teardownTest := test.SetupTest(t)
	defer teardownTest(t)

	user, err := test.CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	if user == nil {
		t.Errorf("got error %v want user", err)
	}

	recovery, err := test.ForgotUserPasswordService.Execute("test@email.com")

	if err != nil {
		t.Errorf("got error %v want recovery", err)
	}

	err = test.ResetUserPasswordService.Execute(recovery.Id, "TestPassword@123")

	if err != nil {
		t.Errorf("got error %v want a password recovery", err)
	}

	_, err = test.GetUserService.Execute("test@email.com", "TestPassword@123")

	if err != nil {
		t.Errorf("got error %v want user with new password", err)
	}
}
