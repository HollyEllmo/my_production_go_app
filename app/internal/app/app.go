package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	_ "github.com/HollyEllmo/my-first-go-project/docs"
	"github.com/HollyEllmo/my-first-go-project/internal/config"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/storage"
	"github.com/HollyEllmo/my-first-go-project/pkg/client/postgresql"
	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	"github.com/HollyEllmo/my-first-go-project/pkg/metric"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	cfg *config.Config
	logger logging.Logger
	router *httprouter.Router
	httpServer *http.Server
	pgClient postgresql.Client
}

func NewApp(config *config.Config, logger logging.Logger) (App, error) {
	logger.Info("Initializing router...")
	router := httprouter.New()

	logger.Info("swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	logger.Println("heartbeat metric initializing")
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	pgConfig := postgresql.NewPgConfig(
		config.PostgreSQL.Username, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)

	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		logger.Fatal(err)
	}

	productStorage := storage.NewProductStorage(pgClient, &logger)
	all, err := productStorage.All(context.Background())
	if err != nil {
		logger.Error(err)
	}
   logger.Fatal(all)

	return App{
		cfg: config,
		logger: logger,
		router: router,
		pgClient: pgClient,
	}, nil
}

func (a *App) Run() {
	a.StartHTTP()
}


func (a *App) StartHTTP() error {
	a.logger.Info("start HTTP")

	var listener net.Listener

	 if a.cfg.Listen.Type == config.LISTEN_TYPE_SOCK {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			a.logger.Fatal(err)
		}
		socketPath := path.Join(appDir, a.cfg.Listen.SocketFile)
		a.logger.Infof("socket path: %s", socketPath)

		a.logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			a.logger.Fatal(err)
		}
	} else {
		a.logger.Infof("bind application to host: %s and port: %d", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		Debug:            a.cfg.IsDebug,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler: handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	a.logger.Print("application completly initialized and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warning("server shutdown")
		default:
			a.logger.WithError(err).Fatal("failed to start server")
		}
	}

	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		a.logger.WithError(err).Fatal("failed to shutdown server")
	}

	return err
}


