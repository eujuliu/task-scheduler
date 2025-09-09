package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"scheduler/internal/config"
	http_handlers "scheduler/internal/handlers/http"
	http_webhooks "scheduler/internal/handlers/http/webhooks"
	"scheduler/pkg/http/helpers"
	"scheduler/pkg/http/middlewares"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	router *gin.Engine
	server *http.Server
	config *config.ServerConfig
}

func New(config *config.ServerConfig) *Server {
	gin.SetMode(config.GinMode)

	router := gin.New()

	router.Use(middlewares.SecureHeaders)
	router.Use(middlewares.Logger)
	router.Use(middlewares.Errors)
	router.Use(middlewares.Cors)
	router.Use(gin.Recovery())

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
		config: config,
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
		routes.POST("/auth/register", http_handlers.Register)
		routes.POST("/auth/login", http_handlers.Login)
		routes.POST("/auth/forgot-password", http_handlers.ForgotPassword)
		routes.POST("/auth/reset-password", http_handlers.ResetUserPassword)
		routes.Any("/stripe-webhook", http_webhooks.StripePayments)
		routes.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})
	}

	protected := routes.Group("/")
	protected.Use(middlewares.Authentication)
	{
		protected.POST("/refresh", middlewares.VerifyRefreshToken, http_handlers.Refresh)

		protected.POST("/buy-credits", http_handlers.BuyCredits)
		protected.GET("/transactions", http_handlers.GetTransactions)
		protected.GET("/transaction/:id", http_handlers.GetTransaction)

		protected.POST("/task", http_handlers.CreateTask)
		protected.PUT("/task/:id", http_handlers.UpdateTask)
		protected.PUT("/task/cancel/:id", http_handlers.CancelTask)
		protected.GET("/task/:id", http_handlers.GetTask)
		protected.GET("/tasks", http_handlers.GetTasks)
	}
}

func (s *Server) setupValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		slog.Info("setup validators...")
		_ = v.RegisterValidation("date", helpers.Datatime)
	}
}
