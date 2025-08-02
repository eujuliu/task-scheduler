package main

import (
	"scheduler/test"
	"testing"
)

func TestGetUser(t *testing.T) {
	teardownTest := test.SetupTest(t)
	defer teardownTest(t)

	_, err := test.CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	if err != nil {
		t.Errorf("got error %v want success", err)
	}

	_, err = test.GetUserService.Execute("test@email.com", "Password@123")

	if err != nil {
		t.Errorf("got error %v want user", err)
	}
}

func TestGetUseWithWrongPassword(t *testing.T) {
	teardownTest := test.SetupTest(t)
	defer teardownTest(t)

	_, err := test.CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	if err != nil {
		t.Errorf("got error %v want success", err)
	}

	user, err := test.GetUserService.Execute("test@email.com", "Password@12")

	if user != nil {
		t.Errorf("got error %v want user", err)
	}
}
