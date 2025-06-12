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
	router *httprouter.Router
	httpServer *http.Server
	pgClient postgresql.Client
}

func NewApp(ctx context.Context, config *config.Config) (App, error) {
	logging.GetLogger(ctx).Println("Initializing router...")
	router := httprouter.New()

	logging.GetLogger(ctx).Println("swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	logging.GetLogger(ctx).Println("heartbeat metric initializing")
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	pgConfig := postgresql.NewPgConfig(
		config.PostgreSQL.Username, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)

	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		logging.GetLogger(ctx).Fatalln(err)
	}

	productStorage := storage.NewProductStorage(pgClient)
	all, err := productStorage.All(context.Background())
	if err != nil {
		logging.GetLogger(ctx).Fatalln(err)
	} else {
		logging.GetLogger(ctx).Infof("Successfully connected to database, found %d products", len(all))
	}

	return App{
		cfg: config,
		router: router,
		pgClient: pgClient,
	}, nil
}

func (a *App) Run(ctx context.Context) {
	a.StartHTTP(ctx)
}


func (a *App) StartHTTP(ctx context.Context) error {
    logging.GetLogger(ctx).Infoln("start HTTP")

	var listener net.Listener

	 if a.cfg.Listen.Type == config.LISTEN_TYPE_SOCK {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logging.GetLogger(ctx).Fatalln(err)
		}
		socketPath := path.Join(appDir, a.cfg.Listen.SocketFile)
		logging.GetLogger(ctx).Infof("socket path: %s", socketPath)

		logging.GetLogger(ctx).Infoln("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			logging.GetLogger(ctx).Fatalln(err)
		}
	} else {
		logging.GetLogger(ctx).Infof("bind application to host: %s and port: %d", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		if err != nil {
			logging.GetLogger(ctx).Fatalln(err)
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

	logging.GetLogger(ctx).Println("application completly initialized and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logging.GetLogger(ctx).Warnln("server shutdown")
		default:
			logging.GetLogger(ctx).WithError(err).Fatalln("failed to start server")
		}
	}

	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Fatalln("failed to shutdown server")
	}

	return err
}


