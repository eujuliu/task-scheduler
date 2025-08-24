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
var ErrorRepository repos.IErrorRepository
var TaskRepository repos.ITaskRepository

var CreateUserService *services.CreateUserService
var GetUserService *services.GetUserService

var ForgotUserPasswordService *services.ForgotUserPasswordService
var ResetUserPasswordService *services.ResetUserPasswordService

var CreateTransactionService *services.CreateTransactionService
var UpdateTransactionService *services.UpdateTransactionService

var CreateTaskService *services.CreateTaskService

func teardown(tb testing.TB) {
	UserRepository = in_memory_repos.NewInMemoryUserRepository()
	PasswordRepository = in_memory_repos.NewInMemoryPasswordRepository()
	TransactionRepository = in_memory_repos.NewInMemoryTransactionRepository()
	ErrorRepository = in_memory_repos.NewInMemoryErrorRepository()
	TaskRepository = in_memory_repos.NewInMemoryTaskRepository()
}

func Setup(tb testing.TB) func(tb testing.TB) {
	UserRepository = in_memory_repos.NewInMemoryUserRepository()
	PasswordRepository = in_memory_repos.NewInMemoryPasswordRepository()
	TransactionRepository = in_memory_repos.NewInMemoryTransactionRepository()
	ErrorRepository = in_memory_repos.NewInMemoryErrorRepository()
	TaskRepository = in_memory_repos.NewInMemoryTaskRepository()

	CreateUserService = services.NewCreateUserService(UserRepository)
	GetUserService = services.NewGetUserService(UserRepository)

	ForgotUserPasswordService = services.NewForgotUserPasswordService(UserRepository, PasswordRepository)
	ResetUserPasswordService = services.NewResetUserPasswordService(UserRepository, PasswordRepository)

	CreateTransactionService = services.NewCreateTransactionService(UserRepository, TransactionRepository)
	UpdateTransactionService = services.NewUpdateTransactionService(UserRepository, TransactionRepository, ErrorRepository)

	CreateTaskService = services.NewCreateTaskService(UserRepository, TransactionRepository, TaskRepository)

	return teardown
}
