package test

import (
	"scheduler/internal/interfaces"
	in_memory_paymentgateway "scheduler/internal/payment_gateway/in_memory"
	"scheduler/internal/queue"
	in_memory_repos "scheduler/internal/repositories/in_memory"
	"scheduler/internal/services"
	"scheduler/pkg/scheduler"
	"sync"
	"testing"

	"github.com/jonboulle/clockwork"
)

var once sync.Once

var (
	UserRepository        interfaces.IUserRepository
	PasswordRepository    interfaces.IPasswordRecoveryRepository
	TransactionRepository interfaces.ITransactionRepository
	ErrorRepository       interfaces.IErrorRepository
	TaskRepository        interfaces.ITaskRepository
)

var (
	CustomerPaymentGateway interfaces.ICustomerPaymentGateway
	PaymentPaymentGateway  interfaces.IPaymentPaymentGateway
)

var (
	Queue     interfaces.IQueue
	Scheduler *scheduler.Scheduler
	Clock     *clockwork.FakeClock
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
	GetTransactionsService           *services.GetTransactionsService
	GetTransactionService            *services.GetTransactionService
)

var (
	CreateTaskService       *services.CreateTaskService
	UpdateTaskService       *services.UpdateTaskService
	GetTasksByUserIdService *services.GetTasksByUserIdService
	GetTasksByRunAtService  *services.GetTasksByRunAtService
	GetTaskService          *services.GetTaskService
)

func teardown(tb testing.TB) {
	UserRepository.(*in_memory_repos.InMemoryUserRepository).Clear()
	PasswordRepository.(*in_memory_repos.InMemoryPasswordRecoveryRepository).Clear()
	TransactionRepository.(*in_memory_repos.InMemoryTransactionRepository).Clear()
	ErrorRepository.(*in_memory_repos.InMemoryErrorRepository).Clear()
	TaskRepository.(*in_memory_repos.InMemoryTaskRepository).Clear()
}

func Setup(tb testing.TB) func(tb testing.TB) {
	once.Do(func() {
		UserRepository = in_memory_repos.NewInMemoryUserRepository()
		PasswordRepository = in_memory_repos.NewInMemoryPasswordRepository()
		TransactionRepository = in_memory_repos.NewInMemoryTransactionRepository()
		ErrorRepository = in_memory_repos.NewInMemoryErrorRepository()
		TaskRepository = in_memory_repos.NewInMemoryTaskRepository()

		CustomerPaymentGateway = in_memory_paymentgateway.NewInMemoryCustomerPaymentGateway()
		PaymentPaymentGateway = in_memory_paymentgateway.NewInMemoryPaymentPaymentGateway()

		Clock = clockwork.NewFakeClock()
		Queue = queue.NewInMemoryQueue()
		Scheduler = scheduler.NewScheduler(Clock, Queue, 20, TaskRepository)
		go Scheduler.Run()

		CreateUserService = services.NewCreateUserService(UserRepository, CustomerPaymentGateway)
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
			CustomerPaymentGateway,
			PaymentPaymentGateway,
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
		GetTransactionsService = services.NewGetTransactionsService(
			UserRepository,
			TransactionRepository,
		)
		GetTransactionService = services.NewGetTransactionService(
			UserRepository,
			TransactionRepository,
		)

		CreateTaskService = services.NewCreateTaskService(
			UserRepository,
			TaskRepository,
			CreateTransactionService,
			UpdateTaskTransactionService,
			Scheduler,
		)
		UpdateTaskService = services.NewUpdateTaskService(
			TaskRepository,
			TransactionRepository,
			UpdateTaskTransactionService,
			Scheduler,
		)
		GetTasksByUserIdService = services.NewGetTasksByUserIdService(
			UserRepository,
			TaskRepository,
		)
		GetTasksByRunAtService = services.NewGetTasksByRunAtService(
			TaskRepository,
		)
		GetTaskService = services.NewGetTaskService(
			UserRepository,
			TaskRepository,
		)
	})

	return teardown
}
