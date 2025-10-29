package main

import (
	"context"
	"log/slog"
	"os"
	"scheduler/internal/composition"
	"scheduler/pkg/http"
	"scheduler/pkg/tracing"

	"github.com/gin-gonic/gin"
)

func main() {
	_, err := tracing.InitTracer()
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	deps, err := composition.Initialize()
	if err != nil {
		panic(err)
	}

	server := http.New(deps)

	if deps.Config.Server.GinMode == gin.DebugMode {
		deps.DB.SeedForTest()
	}

	ctx := context.Background()

	go func() {
		err = deps.UpdateTaskQueueHandler.Handle(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	go deps.Scheduler.Run()

	server.Start()
}
