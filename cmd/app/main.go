package main

import (
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

	go deps.Scheduler.Run()
	server.Start()
}
