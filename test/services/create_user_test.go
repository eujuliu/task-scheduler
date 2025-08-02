package main

import (
	"scheduler/test"
	"testing"
)

func TestCreateUser(t *testing.T) {
	teardownTest := test.SetupTest(t)
	defer teardownTest(t)

	_, err := test.CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	if err != nil {
		t.Errorf("got %v want success", err)
	}
}

func TestDuplicity(t *testing.T) {
	teardownTest := test.SetupTest(t)
	defer teardownTest(t)

	_, err := test.CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	if err != nil {
		t.Errorf("got error %v want success", err)
	}

	_, err = test.CreateUserService.Execute("testuser", "test@email.com", "Password@123")

	if err == nil {
		t.Errorf("expected to be user already exists error, got user creation")
	}
}
