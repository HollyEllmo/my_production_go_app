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

	

	logging.Infoln(ctx, "config initializing")
	cfg := config.GetConfig()


   
	ctx = logging.ContextWithLogger(ctx, logging.NewLogger())

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		logging.Fatalln(ctx, err)
	}

	logging.Infoln(ctx, "Running Application")
	a.Run(ctx)
}