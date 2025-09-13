package composition

import (
	"log/slog"
	"os"
	"scheduler/internal/config"
	http_handlers "scheduler/internal/handlers/http"
	http_webhooks "scheduler/internal/handlers/http/webhooks"
	"scheduler/internal/interfaces"
	stripe_paymentgateway "scheduler/internal/payment_gateway/stripe"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/postgres"
	"scheduler/pkg/rabbitmq"
	"scheduler/pkg/scheduler"

	"github.com/gin-gonic/gin"
	"github.com/jonboulle/clockwork"
)

type Dependencies struct {
	DB                        *postgres.Database
	Scheduler                 *scheduler.Scheduler
	Config                    *config.Config
	Queue                     *interfaces.IQueue
	BuyCreditsHandler         *http_handlers.BuyCreditsHandler
	CancelTaskHandler         *http_handlers.CancelTaskHandler
	CreateTaskHandler         *http_handlers.CreateTaskHandler
	ForgotUserPasswordHandler *http_handlers.ForgotPasswordHandler
	GetTaskHandler            *http_handlers.GetTaskHandler
	GetTasksHandler           *http_handlers.GetTasksHandler
	GetTransactionHandler     *http_handlers.GetTransactionHandler
	GetTransactionsHandler    *http_handlers.GetTransactionsHandler
	LoginHandler              *http_handlers.LoginHandler
	RefreshTokenHandler       *http_handlers.RefreshTokenHandler
	RegisterUserHandler       *http_handlers.RegisterHandler
	ResetUserPasswordHandler  *http_handlers.ResetUserPasswordHandler
	UpdateTaskHandler         *http_handlers.UpdateTaskHandler

	StripePaymentUpdateWebhook *http_webhooks.StripePaymentUpdateWebhook
}

func Initialize() (*Dependencies, error) {
	config := config.NewConfig()
	db, err := postgres.NewPostgres(config.Database)
	if err != nil {
		return nil, err
	}

	rmq, err := rabbitmq.NewRabbitMQ(config.RabbitMQ)
	if err != nil {
		return nil, err
	}

	err = rmq.AddDurableQueue("tasks-queue", "task-exchange", "task.send")
	if err != nil {
		return nil, err
	}

	loggerLevel := slog.LevelInfo

	if config.Server.GinMode == gin.DebugMode {
		loggerLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
	}))

	slog.SetDefault(logger)

	scheduler := scheduler.NewScheduler(clockwork.NewRealClock(), rmq, 20)

	userRepository := postgres_repos.NewPostgresUserRepository(db)
	passwordRepository := postgres_repos.NewPostgresPasswordRepository(db)
	transactionRepository := postgres_repos.NewPostgresTransactionRepository(db)
	errorRepository := postgres_repos.NewPostgresErrorRepository(db)
	taskRepository := postgres_repos.NewPostgresTaskRepository(db)

	customerPaymentGateway := stripe_paymentgateway.NewStripeCustomerPaymentGateway()
	paymentPaymentGateway := stripe_paymentgateway.NewStripePaymentPaymentGateway()

	createUserService := services.NewCreateUserService(userRepository, customerPaymentGateway)
	getUserService := services.NewGetUserService(userRepository)

	forgotUserPasswordService := services.NewForgotUserPasswordService(
		userRepository,
		passwordRepository,
	)
	resetUserPasswordService := services.NewResetUserPasswordService(
		userRepository,
		passwordRepository,
	)

	createTransactionService := services.NewCreateTransactionService(
		userRepository,
		transactionRepository,
		customerPaymentGateway,
		paymentPaymentGateway,
	)

	updatePurchaseTransactionService := services.NewUpdatePurchaseTransactionService(
		userRepository,
		transactionRepository,
		errorRepository,
	)

	updateTaskTransactionService := services.NewUpdateTaskTransactionService(
		userRepository,
		transactionRepository,
		errorRepository,
	)
	getTransactionsService := services.NewGetTransactionsService(
		userRepository,
		transactionRepository,
	)
	getTransactionService := services.NewGetTransactionService(
		userRepository,
		transactionRepository,
	)

	createTaskService := services.NewCreateTaskService(
		userRepository,
		taskRepository,
		createTransactionService,
		updateTaskTransactionService,
	)
	updateTaskService := services.NewUpdateTaskService(
		taskRepository,
		transactionRepository,
		updateTaskTransactionService,
	)
	getTasksByUserIdService := services.NewGetTasksByUserIdService(
		userRepository,
		taskRepository,
	)
	getTaskService := services.NewGetTaskService(
		userRepository,
		taskRepository,
	)

	buyCreditsHandler := http_handlers.NewBuyCreditsHandler(db, createTransactionService)
	cancelTaskHandler := http_handlers.NewCancelTaskHandler(db, updateTaskService)
	createTaskHandler := http_handlers.NewCreateTaskHandler(db, scheduler, createTaskService)
	forgotUserPasswordHandler := http_handlers.NewForgotPasswordHandler(forgotUserPasswordService)
	getTaskHandler := http_handlers.NewGetTaskHandler(getTaskService)
	getTasksHandler := http_handlers.NewGetTasksHandler(getTasksByUserIdService)
	getTransactionHandler := http_handlers.NewGetTransactionHandler(getTransactionService)
	getTransactionsHandler := http_handlers.NewGetTransactionsHandler(getTransactionsService)
	loginHandler := http_handlers.NewLoginHandler(config, getUserService)
	refreshTokenHandler := http_handlers.NewRefreshTokenHandler(config)
	registerUserHandler := http_handlers.NewRegisterHandler(config, createUserService)
	resetUserPasswordHandler := http_handlers.NewResetUserPasswordHandler(
		db,
		resetUserPasswordService,
	)
	updateTaskHandler := http_handlers.NewUpdateTaskHandler(db, scheduler, updateTaskService)

	stripePaymentUpdateWebhook := http_webhooks.NewStripePaymentUpdateWebhook(
		config.Stripe,
		updatePurchaseTransactionService,
	)

	return &Dependencies{
		DB:                        db,
		Scheduler:                 scheduler,
		Config:                    config,
		BuyCreditsHandler:         buyCreditsHandler,
		CancelTaskHandler:         cancelTaskHandler,
		CreateTaskHandler:         createTaskHandler,
		ForgotUserPasswordHandler: forgotUserPasswordHandler,
		GetTaskHandler:            getTaskHandler,
		GetTasksHandler:           getTasksHandler,
		GetTransactionHandler:     getTransactionHandler,
		GetTransactionsHandler:    getTransactionsHandler,
		LoginHandler:              loginHandler,
		RefreshTokenHandler:       refreshTokenHandler,
		RegisterUserHandler:       registerUserHandler,
		ResetUserPasswordHandler:  resetUserPasswordHandler,
		UpdateTaskHandler:         updateTaskHandler,

		StripePaymentUpdateWebhook: stripePaymentUpdateWebhook,
	}, nil
}
