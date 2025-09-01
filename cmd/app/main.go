package main

import (
	"log/slog"
	"os"
	"scheduler/internal/config"
	"scheduler/pkg/http"
	"scheduler/pkg/postgres"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	config := config.Load()

	server := http.New(config.Server)
	_, err := postgres.Load(config.Database)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	server.Start()
}
