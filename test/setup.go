package test

import (
	repos "scheduler/internal/repositories"
	in_memory_repos "scheduler/internal/repositories/in_memory"
	"scheduler/internal/services"
	"testing"
)

var UserRepository repos.IUserRepository
var PasswordRepository repos.IPasswordRecoveryRepository

var CreateUserService *services.CreateUserService
var GetUserService *services.GetUserService

var ForgotUserPasswordService *services.ForgotUserPasswordService
var ResetUserPasswordService *services.ResetUserPasswordService

func teardown(tb testing.TB) {
	users := UserRepository.Get()

	for _, u := range users {
		_ = UserRepository.Delete(u.GetId())
	}

	tokens := PasswordRepository.Get()

	for _, t := range tokens {
		_ = PasswordRepository.Delete(t.GetId())
	}
}

func Setup(tb testing.TB) func(tb testing.TB) {
	UserRepository = in_memory_repos.NewInMemoryUserRepository()
	PasswordRepository = in_memory_repos.NewInMemoryPasswordRepository()

	CreateUserService = services.NewCreateUserService(UserRepository)
	GetUserService = services.NewGetUserService(UserRepository)

	ForgotUserPasswordService = services.NewForgotUserPasswordService(UserRepository, PasswordRepository)
	ResetUserPasswordService = services.NewResetUserPasswordService(UserRepository, PasswordRepository)

	return teardown
}
