package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/HollyEllmo/my-first-go-project/docs"
	"github.com/HollyEllmo/my-first-go-project/internal/config"
	"github.com/HollyEllmo/my-first-go-project/internal/controller/grpc/v1/product"
	"github.com/HollyEllmo/my-first-go-project/pkg/client/postgresql"
	"github.com/HollyEllmo/my-first-go-project/pkg/logging"
	"github.com/HollyEllmo/my-first-go-project/pkg/metric"
	pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
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

	productServiceServer pb_prod_products.ProductServiceServer
}

func NewApp(ctx context.Context, config *config.Config) (App, error) {
	logging.Infoln(ctx, "Initializing router...")
	router := httprouter.New()

	logging.Infoln(ctx, "swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	logging.Infoln(ctx, "heartbeat metric initializing")
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	pgConfig := postgresql.NewPgConfig(
		config.PostgreSQL.Username, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)

	pgClient, err := postgresql.NewClient(ctx, 5, time.Second*5, pgConfig)
	if err != nil {
		logging.WithError(ctx, err).Fatalln("failed to connect to PostgreSQL")
	}

	// productStorage := storage.NewProductStorage(pgClient)
	// all, err := productStorage.All(ctx, nil, nil)
	// if err != nil {
	// 	logging.GetLogger().Fatalln(err)
	// } else {
	// 	logging.Infof(ctx, "Successfully connected to database, found %d products", len(all))
	// }

	productServiceServer := product.NewServer(
		pb_prod_products.UnimplementedProductServiceServer{},
	)

	return App{
		cfg: config,
		router: router,
		pgClient: pgClient,
		productServiceServer: productServiceServer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		return a.StartHTTP(ctx)
	})
	grp.Go(func() error {
		return a.StartGRPC(ctx, a.productServiceServer)
	})
	return grp.Wait()
}

func (a *App) StartGRPC(ctx context.Context, server pb_prod_products.ProductServiceServer) error {
	logger := logging.WithFields(ctx, map[string]interface{}{
		"IP":   a.cfg.GRPC.IP,
		"Port": a.cfg.GRPC.Port,
	})

	logger.Println("gRPC server initializing")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.GRPC.IP, a.cfg.GRPC.Port))
	if err != nil {
		logger.WithError(err).Fatalln("failed to listen on port")
	}

	serverOptions := []grpc.ServerOption{}
	a.grpcServer = grpc.NewServer(serverOptions...)

	pb_prod_products.RegisterProductServiceServer(a.grpcServer, server)

	reflection.Register(a.grpcServer)

	return a.grpcServer.Serve(listener)
}

func (a *App) StartHTTP(ctx context.Context) error {
    logger := logging.WithFields(ctx, map[string]interface{}{
		"IP":   a.cfg.HTTP.IP,
		"Port": a.cfg.HTTP.Port,
	})

	logger.Println("HTTP server initializing")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.HTTP.IP, a.cfg.HTTP.Port))
	if err != nil {
		logger.WithError(err).Fatalln("failed to listen on port")
	}

	logger.WithFields(map[string]interface{}{
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

	logger.Println("application completly initialized and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warnln("server shutdown")
		default:
			logger.WithError(err).Fatalln("failed to start server")
		}
	}

	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		logger.WithError(err).Fatalln("failed to shutdown server")
	}

	return err
}


