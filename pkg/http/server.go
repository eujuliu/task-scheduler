package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"scheduler/internal/composition"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/http/middlewares"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
	server *http.Server
	deps   *composition.Dependencies
}

func New(deps *composition.Dependencies) *Server {
	config := deps.Config.Server
	gin.SetMode(config.GinMode)

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.SecureHeaders)
	router.Use(middlewares.Logger)
	router.Use(middlewares.Errors)
	router.Use(middlewares.Cors)

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.Use(middlewares.RateLimiter(deps.RateLimiter))

	server := http.Server{
		Addr:           fmt.Sprintf("%s:%s", config.Host, config.Port),
		Handler:        router,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		IdleTimeout:    config.IdleTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	return &Server{
		router: router,
		server: &server,
		deps:   deps,
	}
}

func (s *Server) Start() {
	s.setupValidators()
	s.setupRoutes()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(fmt.Sprintf("server listen: %s", err))
			panic(err)
		}
	}()

	<-ctx.Done()

	stop()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		slog.Error(fmt.Sprintf("server forced to shutdown: %s", err))
		panic(err)
	}

	slog.Info("server exiting")
}

func (s *Server) setupRoutes() {
	routes := s.router.Group("/")
	{
		routes.POST("/auth/register", s.deps.RegisterUserHandler.Handle)
		routes.POST("/auth/login", s.deps.LoginHandler.Handle)
		routes.POST("/auth/forgot-password", s.deps.ForgotUserPasswordHandler.Handle)
		routes.POST("/auth/reset-password", s.deps.ResetUserPasswordHandler.Handle)
		routes.POST("/stripe-webhook", s.deps.StripePaymentUpdateWebhook.Hook)
		routes.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	protected := routes.Group("/")
	protected.Use(middlewares.Authentication)
	{
		protected.POST(
			"/refresh",
			middlewares.VerifyRefreshToken,
			s.deps.RefreshTokenHandler.Handle,
		)
		protected.DELETE("/logoff", s.deps.LogoffHandler.Handle)

		protected.POST("/buy-credits", s.deps.BuyCreditsHandler.Handle)
		protected.GET("/transactions", s.deps.GetTransactionsHandler.Handle)
		protected.GET("/transaction/:id", s.deps.GetTransactionHandler.Handle)

		protected.POST("/task", s.deps.CreateTaskHandler.Handle)
		protected.PUT("/task/:id", s.deps.UpdateTaskHttpHandler.Handle)
		protected.PUT("/task/cancel/:id", s.deps.CancelTaskHandler.Handle)
		protected.GET("/task/:id", s.deps.GetTaskHandler.Handle)
		protected.GET("/tasks", s.deps.GetTasksHandler.Handle)
	}
}

func (s *Server) setupValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("date", helpers.Datatime)
		_ = v.RegisterValidation("utc", helpers.UTCDateTime)
	}
}
