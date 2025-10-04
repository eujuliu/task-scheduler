package main

import (
	"context"
	"log/slog"
	"scheduler/docs"
	"scheduler/internal/composition"
	"scheduler/pkg/http"

	"github.com/gin-gonic/gin"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter 'Bearer ' followed by your access token

// @securityDefinitions.apikey RefreshTokenAuth
// @in cookie
// @name refresh_token
// @description Refresh token stored in HTTP-only cookie

func main() {
	docs.SwaggerInfo.Title = "Task Scheduler API"
	docs.SwaggerInfo.Description = "A task scheduler API with user authentication, task management, and credit purchases."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

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
