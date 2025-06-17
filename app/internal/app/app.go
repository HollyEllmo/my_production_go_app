package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"

	_ "github.com/HollyEllmo/my-first-go-project/docs"
	"github.com/HollyEllmo/my-first-go-project/internal/config"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/storage"
	"github.com/HollyEllmo/my-first-go-project/pkg/client/postgresql"
	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	"github.com/HollyEllmo/my-first-go-project/pkg/metric"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg *config.Config
	router *httprouter.Router
	httpServer *http.Server
	grpcServer *grpc.Server
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

	pgClient, err := postgresql.NewClient(ctx, 5, time.Second*5, pgConfig)
	if err != nil {
		logging.GetLogger(ctx).Fatalln(err)
	}

	productStorage := storage.NewProductStorage(pgClient)
	all, err := productStorage.All(ctx)
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

func (a *App) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		return a.StartHTTP(ctx)
	})
	logging.GetLogger(ctx).Infoln("application initialized and started")
	return grp.Wait()
}

func (a *App) StartGRPC(ctx context.Context) error {
	logging.GetLogger(ctx).WithFields(map[string]interface{}{
		"IP":   a.cfg.GRPC.IP,
		"Port": a.cfg.GRPC.Port,
	})

	// listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.GRPC.IP, a.cfg.GRPC.Port))
	// if err != nil {
	// 	logging.GetLogger(ctx).WithError(err).Fatalln("failed to listen on port")
	// }

	serverOptions := []grpc.ServerOption{}
	a.grpcServer = grpc.NewServer(serverOptions...)

	return nil
}

func (a *App) StartHTTP(ctx context.Context) error {
    logging.GetLogger(ctx).WithFields(map[string]interface{}{
		"IP":   a.cfg.HTTP.IP,
		"Port": a.cfg.HTTP.Port,
	})

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.HTTP.IP, a.cfg.HTTP.Port))
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Fatalln("failed to listen on port")
	}

	logging.GetLogger(ctx).WithFields(map[string]interface{}{
		"AllowedMethods":     a.cfg.HTTP.CORS.AllowedMethods,
		"AllowedOrigins":     a.cfg.HTTP.CORS.AllowedOrigins,
		"AllowCredentials":   a.cfg.HTTP.CORS.AllowCredentials,
		"AllowedHeaders":     a.cfg.HTTP.CORS.AllowedHeaders,
		"OptionsPassthrough": a.cfg.HTTP.CORS.OptionsPassthrough,
		"ExposedHeaders":     a.cfg.HTTP.CORS.ExposedHeaders,
		"Debug":              a.cfg.HTTP.CORS.Debug,
	}).Info("CORS configuration")

	c := cors.New(cors.Options{
		AllowedMethods:     a.cfg.HTTP.CORS.AllowedMethods,
		AllowedOrigins:     a.cfg.HTTP.CORS.AllowedOrigins,
		AllowCredentials:   a.cfg.HTTP.CORS.AllowCredentials,
		AllowedHeaders:     a.cfg.HTTP.CORS.AllowedHeaders,
		OptionsPassthrough: a.cfg.HTTP.CORS.OptionsPassthrough,
		ExposedHeaders:     a.cfg.HTTP.CORS.ExposedHeaders,
		Debug:              a.cfg.HTTP.CORS.Debug,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler: handler,
		WriteTimeout: a.cfg.HTTP.WriteTimeout,
		ReadTimeout: a.cfg.HTTP.ReadTimeout,
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

	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Fatalln("failed to shutdown server")
	}

	return err
}


