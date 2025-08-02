package main

import (
	"scheduler/test"
	"testing"
	"time"
)

func TestForgotPassword(t *testing.T) {
	teardownTest := test.SetupTest(t)
	defer teardownTest(t)

	user, _ := test.CreateUserService.Execute("testuser", "user@email.com", "Password@123")

	recovery, err := test.ForgotUserPasswordService.Execute(user.Email)

	if err != nil {
		t.Errorf("got error %v want recovery token", err)
	}

	if !recovery.InTime(time.Now()) {
		t.Error("got an invalid expiration time")
	}
}

func TestExpiredToken(t *testing.T) {
	teardownTest := test.SetupTest(t)
	defer teardownTest(t)

	user, _ := test.CreateUserService.Execute("testuser", "user@email.com", "Password@123")

	recovery, err := test.ForgotUserPasswordService.Execute(user.Email)

	if err != nil {
		t.Errorf("got an error %v expect a token", err)
	}

	if recovery.InTime(time.Now().Add(6 * time.Minute)) {
		t.Errorf("got in time token after 6 minutes")
	}
}
