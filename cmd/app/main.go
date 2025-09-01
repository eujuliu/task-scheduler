package main

import (
	"log/slog"
	"os"
	"scheduler/internal/config"
	"scheduler/pkg/http"
	"scheduler/pkg/postgres"
)

func main() {
	config := config.Load()

	loggerLevel := slog.LevelInfo

	if config.Server.GinMode == "debug" {
		loggerLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
	}))

	slog.SetDefault(logger)

	server := http.New(config.Server)
	_, err := postgres.Load(config.Database)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	server.Start()
}
