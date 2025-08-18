package test

import (
	repos "scheduler/internal/repositories"
	in_memory_repos "scheduler/internal/repositories/in_memory"
	"scheduler/internal/services"
	"testing"
)

var UserRepository repos.IUserRepository
var PasswordRepository repos.IPasswordRecoveryRepository
var TransactionRepository repos.ITransactionRepository

var CreateUserService *services.CreateUserService
var GetUserService *services.GetUserService

var ForgotUserPasswordService *services.ForgotUserPasswordService
var ResetUserPasswordService *services.ResetUserPasswordService

var CreateTransactionService *services.CreateTransactionService

// var CreateTaskService *services.CreateTaskService

func teardown(tb testing.TB) {
	users := UserRepository.Get()
	tokens := PasswordRepository.Get()
	// transactions := TransactionRepository.Get()

	for _, u := range users {
		_ = UserRepository.Delete(u.GetId())
	}

	for _, t := range tokens {
		_ = PasswordRepository.Delete(t.GetId())
	}

	// for _, t := range transactions {
	// 	_ = TransactionRepository.Delete(t.GetId())
	// }
}

func Setup(tb testing.TB) func(tb testing.TB) {
	UserRepository = in_memory_repos.NewInMemoryUserRepository()
	PasswordRepository = in_memory_repos.NewInMemoryPasswordRepository()
	TransactionRepository = in_memory_repos.NewInMemoryTransactionRepository()

	CreateUserService = services.NewCreateUserService(UserRepository)
	GetUserService = services.NewGetUserService(UserRepository)

	ForgotUserPasswordService = services.NewForgotUserPasswordService(UserRepository, PasswordRepository)
	ResetUserPasswordService = services.NewResetUserPasswordService(UserRepository, PasswordRepository)

	CreateTransactionService = services.NewCreateTransactionService(UserRepository, TransactionRepository)

	// CreateTaskService = services.NewCreateTaskService(UserRepository, TransactionRepository)

	return teardown
}
