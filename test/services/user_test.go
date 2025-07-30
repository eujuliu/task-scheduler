package main

import (
	in_memory_repos "scheduler/internal/repositories/in-memory"
	"scheduler/internal/services"
	"testing"
)

func TestCreation(t *testing.T) {
	repository := in_memory_repos.NewUserRepository()
	service := services.NewUserService(repository)

	_, err := service.CreateUser("testuser", "test@email.com", "Password@123")

	if err != nil {
		t.Errorf("got %v want success", err)
	}
}

func TestCreationDuplicity(t *testing.T) {
	repository := in_memory_repos.NewUserRepository()
	service := services.NewUserService(repository)

	_, err := service.CreateUser("testuser", "test@email.com", "Password@123")

	if err != nil {
		t.Errorf("got error %v want success", err)
	}

	_, err = service.CreateUser("testuser", "test@email.com", "Password@123")

	if err == nil {
		t.Errorf("expected to be user already exists error, got user creation")
	}
}

func TestGet(t *testing.T) {
	repository := in_memory_repos.NewUserRepository()
	service := services.NewUserService(repository)

	_, err := service.CreateUser("testuser", "test@email.com", "Password@123")

	if err != nil {
		t.Errorf("got error %v want success", err)
	}

	_, err = service.GetUser("test@email.com", "Password@123")

	if err != nil {
		t.Errorf("got error %v want user", err)
	}
}

func TestWrongPassword(t *testing.T) {
	repository := in_memory_repos.NewUserRepository()
	service := services.NewUserService(repository)

	_, err := service.CreateUser("testuser", "test@email.com", "Password@123")

	if err != nil {
		t.Errorf("got error %v want success", err)
	}

	user, err := service.GetUser("test@email.com", "Password@12")

	if user != nil {
		t.Errorf("got error %v want user", err)
	}
}
