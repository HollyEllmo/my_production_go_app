package main

import (
	"context"

	"github.com/HollyEllmo/my-first-go-project/internal/app"
	"github.com/HollyEllmo/my-first-go-project/internal/config"
	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logging.GetLogger(ctx)

	logger.Println("config initializing")
	cfg := config.GetConfig()


   
	ctx = logging.ContextWithLogger(ctx, logger)

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		logging.GetLogger(ctx).Fatalln(err)
	}

	logger.Infoln("Running Application")
	a.Run(ctx)
}