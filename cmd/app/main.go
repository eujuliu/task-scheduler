package main

import (
	"context"
	"log/slog"
	"scheduler/internal/composition"
	"scheduler/pkg/http"

	"github.com/gin-gonic/gin"
)

func main() {
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
