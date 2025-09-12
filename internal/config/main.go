package config

import (
	"scheduler/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServerConfig struct {
	Host            string
	Port            string
	GinMode         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	MaxHeaderBytes  int
}

type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
}

type DatabaseConfig struct {
	Port     string
	User     string
	Password string
	DBName   string
	Host     string
}

type StripeConfig struct {
	APIKey         string
	EndpointSecret string
}

type Config struct {
	Server   *ServerConfig
	JWT      *JWTConfig
	Database *DatabaseConfig
	Stripe   *StripeConfig
}

func NewConfig() *Config {
	return &Config{
		Server: &ServerConfig{
			Host:            utils.GetEnv("Host", "0.0.0.0"),
			Port:            utils.GetEnv("PORT", "8080"),
			GinMode:         utils.GetEnv("GIN_MODE", gin.ReleaseMode),
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			IdleTimeout:     60 * time.Second,
			ShutdownTimeout: 15 * time.Second,
			MaxHeaderBytes:  1 << 20, // 1MB
		},
		JWT: &JWTConfig{
			AccessTokenSecret:  utils.GetEnv("ACCESS_TOKEN_SECRET", uuid.NewString()),
			RefreshTokenSecret: utils.GetEnv("REFRESH_TOKEN_SECRET", uuid.NewString()),
		},
		Database: &DatabaseConfig{
			Port:     utils.GetEnv("POSTGRES_PORT", "5432"),
			User:     utils.GetEnv("POSTGRES_USER", "local_user"),
			Password: utils.GetEnv("POSTGRES_PASSWORD", "local_password"),
			DBName:   utils.GetEnv("POSTGRES_DB", "taskscheduler"),
			Host:     utils.GetEnv("POSTGRES_HOST", "localhost"),
		},
		Stripe: &StripeConfig{
			APIKey:         utils.GetEnv("STRIPE_API_KEY", ""),
			EndpointSecret: utils.GetEnv("STRIPE_WEBHOOK_SIGNING_SECRET", ""),
		},
	}
}
