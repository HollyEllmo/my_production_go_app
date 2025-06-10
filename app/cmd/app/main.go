package main

import (
	"log"

	"github.com/HollyEllmo/my-first-go-project/internal/app"
	"github.com/HollyEllmo/my-first-go-project/internal/config"
	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
)

func main() {
	log.Print("config initializing")
	cfg := config.GetConfig()

	log.Print("logger initializing")
	logging.Init(cfg.AppConfig.LogLevel)
	logger := logging.GetLogger()

	a, err := app.NewApp(cfg, logger)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Println("Running Application")
	a.Run()
}
