package main

import (
	in_memory_repos "scheduler/internal/repositories/in-memory"
	"scheduler/internal/services"
	"testing"
)

func TestUserCreation(t *testing.T) {
	repository := in_memory_repos.NewUserRepository()
	service := services.NewUserService(repository)

	_, err := service.CreateUser("testuser", "test@email.com", "Password@123")

	if err != nil {
		t.Errorf("got %v want success", err)
	}
}

func TestUserAlreadyExist(t *testing.T) {
	repository := in_memory_repos.NewUserRepository()
	service := services.NewUserService(repository)

	_, err := service.CreateUser("testuser", "test@email.com", "Password@123")

	if err != nil {
		t.Errorf("got %v want success", err)
	}

	_, err = service.CreateUser("testuser", "test@email.com", "Password@123")

	if err == nil {
		t.Errorf("expect to be user already exists error, got user creation")
	}
}
