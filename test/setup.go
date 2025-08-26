package test

import (
	"scheduler/internal/services"
	"testing"

	repos "scheduler/internal/repositories"
	in_memory_repos "scheduler/internal/repositories/in_memory"
)

var (
	UserRepository        repos.IUserRepository
	PasswordRepository    repos.IPasswordRecoveryRepository
	TransactionRepository repos.ITransactionRepository
	ErrorRepository       repos.IErrorRepository
	TaskRepository        repos.ITaskRepository
)

var (
	CreateUserService *services.CreateUserService
	GetUserService    *services.GetUserService
)

var (
	ForgotUserPasswordService *services.ForgotUserPasswordService
	ResetUserPasswordService  *services.ResetUserPasswordService
)

var (
	CreateTransactionService         *services.CreateTransactionService
	UpdatePurchaseTransactionService *services.UpdatePurchaseTransactionService
	UpdateTaskTransactionService     *services.UpdateTaskTransactionService
)

var (
	CreateTaskService *services.CreateTaskService
	UpdateTaskService *services.UpdateTaskService
)

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

	ForgotUserPasswordService = services.NewForgotUserPasswordService(
		UserRepository,
		PasswordRepository,
	)
	ResetUserPasswordService = services.NewResetUserPasswordService(
		UserRepository,
		PasswordRepository,
	)

	CreateTransactionService = services.NewCreateTransactionService(
		UserRepository,
		TransactionRepository,
	)
	UpdatePurchaseTransactionService = services.NewUpdatePurchaseTransactionService(
		UserRepository,
		TransactionRepository,
		ErrorRepository,
	)
	UpdateTaskTransactionService = services.NewUpdateTaskTransactionService(
		UserRepository,
		TransactionRepository,
		ErrorRepository,
	)

	CreateTaskService = services.NewCreateTaskService(
		UserRepository,
		TransactionRepository,
		TaskRepository,
		CreateTransactionService,
		UpdateTaskTransactionService,
	)
	UpdateTaskService = services.NewUpdateTaskService(
		TaskRepository,
		TransactionRepository,
		UpdateTaskTransactionService,
	)

	return teardown
}
