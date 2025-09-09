package main

import (
	"log/slog"
	"os"
	"scheduler/internal/config"
	"scheduler/pkg/http"
	"scheduler/pkg/postgres"
	"scheduler/pkg/stripe"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.Load()

	loggerLevel := slog.LevelInfo

	if config.Server.GinMode == gin.DebugMode {
		loggerLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
	}))

	slog.SetDefault(logger)

	_ = stripe.Load(config.Stripe)

	server := http.New(config.Server)
	db, err := postgres.Load(config.Database)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	if config.Server.GinMode == gin.DebugMode {
		db.SeedForTest()
	}

	server.Start()
}
