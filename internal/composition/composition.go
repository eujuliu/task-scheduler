package composition

import (
	"log/slog"
	"os"
	"scheduler/internal/config"
	http_handlers "scheduler/internal/handlers/http"
	http_webhooks "scheduler/internal/handlers/http/webhooks"
	"scheduler/internal/interfaces"
	stripe_paymentgateway "scheduler/internal/payment_gateway/stripe"
	"scheduler/internal/queue"
	postgres_repos "scheduler/internal/repositories/postgres"
	"scheduler/internal/services"
	"scheduler/pkg/postgres"
	"scheduler/pkg/rabbitmq"
	ratelimiter "scheduler/pkg/rate_limiter"
	"scheduler/pkg/redis"
	"scheduler/pkg/scheduler"

	"github.com/gin-gonic/gin"
	"github.com/jonboulle/clockwork"
)

type Dependencies struct {
	DB          *postgres.Database
	Redis       *redis.Redis
	Scheduler   *scheduler.Scheduler
	Config      *config.Config
	Queue       *interfaces.IQueue
	RateLimiter *ratelimiter.SlidingWindowCounterLimiter

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

	err = rmq.AddDurableQueue(queue.SEND_EMAIL_QUEUE, queue.TASK_EXCHANGE, queue.SEND_EMAIL_KEY)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewRedis(config.Redis)
	limiter := ratelimiter.NewSlidingWindowCounterLimiter(
		rdb,
		config.RateLimiter.RequestLimit,
		config.RateLimiter.WindowSize,
		config.RateLimiter.SubWindowSize,
	)

	loggerLevel := slog.LevelInfo

	if config.Server.GinMode == gin.DebugMode {
		loggerLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
	}))

	slog.SetDefault(logger)

	userRepository := postgres_repos.NewPostgresUserRepository(db)
	passwordRepository := postgres_repos.NewPostgresPasswordRepository(db)
	transactionRepository := postgres_repos.NewPostgresTransactionRepository(db)
	errorRepository := postgres_repos.NewPostgresErrorRepository(db)
	taskRepository := postgres_repos.NewPostgresTaskRepository(db)

	customerPaymentGateway := stripe_paymentgateway.NewStripeCustomerPaymentGateway()
	paymentPaymentGateway := stripe_paymentgateway.NewStripePaymentPaymentGateway()

	scheduler := scheduler.NewScheduler(clockwork.NewRealClock(), rmq, 100, taskRepository)

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
		scheduler,
	)
	updateTaskService := services.NewUpdateTaskService(
		taskRepository,
		transactionRepository,
		updateTaskTransactionService,
		scheduler,
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
	createTaskHandler := http_handlers.NewCreateTaskHandler(db, createTaskService)
	forgotUserPasswordHandler := http_handlers.NewForgotPasswordHandler(forgotUserPasswordService)
	getTaskHandler := http_handlers.NewGetTaskHandler(getTaskService)
	getTasksHandler := http_handlers.NewGetTasksHandler(getTasksByUserIdService)
	getTransactionHandler := http_handlers.NewGetTransactionHandler(getTransactionService)
	getTransactionsHandler := http_handlers.NewGetTransactionsHandler(getTransactionsService)
	loginHandler := http_handlers.NewLoginHandler(config, rdb, getUserService)
	refreshTokenHandler := http_handlers.NewRefreshTokenHandler(config, rdb)
	registerUserHandler := http_handlers.NewRegisterHandler(config, rdb, createUserService)
	resetUserPasswordHandler := http_handlers.NewResetUserPasswordHandler(
		db,
		resetUserPasswordService,
	)
	updateTaskHandler := http_handlers.NewUpdateTaskHandler(db, updateTaskService)

	stripePaymentUpdateWebhook := http_webhooks.NewStripePaymentUpdateWebhook(
		config.Stripe,
		updatePurchaseTransactionService,
	)

	return &Dependencies{
		DB:          db,
		Scheduler:   scheduler,
		Config:      config,
		Redis:       rdb,
		RateLimiter: limiter,

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
