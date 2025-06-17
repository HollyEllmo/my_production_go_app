package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/HollyEllmo/my-first-go-project/docs"
	"github.com/HollyEllmo/my-first-go-project/internal/config"
	"github.com/HollyEllmo/my-first-go-project/internal/domain/pruduct/storage"
	"github.com/HollyEllmo/my-first-go-project/internal/handler"
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
	prodServiceClient pb_prod_products.ProductServiceClient
	grpcConn *grpc.ClientConn
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

	// Инициализация обработчика продуктов (будет настроен позже)
	var productHandler *handler.ProductHandler

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

	// Инициализация gRPC клиента для prod_service
	grpcConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Warnln("Failed to connect to prod_service gRPC server")
		// Не делаем fatal, так как сервис может работать без prod_service
	}
	
	var prodServiceClient pb_prod_products.ProductServiceClient
	if grpcConn != nil {
		prodServiceClient = pb_prod_products.NewProductServiceClient(grpcConn)
		logging.GetLogger(ctx).Infoln("Successfully connected to prod_service gRPC")
		
		// Инициализируем обработчик продуктов с gRPC клиентом
		productHandler = handler.NewProductHandler(prodServiceClient)
		productHandler.Register(router)
		logging.GetLogger(ctx).Infoln("Product handler registered")
		
		// Опционально: тестируем соединение
		go func() {
			if err := testProductServiceConnection(ctx, prodServiceClient); err != nil {
				logging.GetLogger(ctx).WithError(err).Warnln("Product service connection test failed")
			}
		}()
	} else {
		logging.GetLogger(ctx).Warnln("Product service not available, product handler not registered")
	}

	return App{
		cfg: config,
		router: router,
		pgClient: pgClient,
		prodServiceClient: prodServiceClient,
		grpcConn: grpcConn,
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

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	logging.GetLogger(ctx).Infoln("Shutting down application...")
	
	// Закрываем gRPC соединение
	if a.grpcConn != nil {
		if err := a.grpcConn.Close(); err != nil {
			logging.GetLogger(ctx).WithError(err).Warnln("Failed to close gRPC connection")
		} else {
			logging.GetLogger(ctx).Infoln("gRPC connection closed successfully")
		}
	}
	
	// Закрываем HTTP сервер
	if a.httpServer != nil {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			logging.GetLogger(ctx).WithError(err).Errorln("Failed to shutdown HTTP server")
			return err
		}
	}
	
	logging.GetLogger(ctx).Infoln("Application shutdown completed")
	return nil
}

// GetProductServiceClient returns the gRPC client for product service
func (a *App) GetProductServiceClient() pb_prod_products.ProductServiceClient {
	return a.prodServiceClient
}

// CallProductService demonstrates how to call methods from prod_service
func (a *App) CallProductService(ctx context.Context) error {
	if a.prodServiceClient == nil {
		return fmt.Errorf("product service client is not initialized")
	}
	
	logging.GetLogger(ctx).Infoln("Calling product service...")
	
	// Пример вызова метода AllProducts
	req := &pb_prod_products.AllProductsRequest{}
	resp, err := a.prodServiceClient.AllProducts(ctx, req)
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Errorln("Failed to call AllProducts")
		return err
	}
	logging.GetLogger(ctx).Infof("Received %d products from service", len(resp.Product))
	
	// Пример вызова метода CreateProduct
	createReq := &pb_prod_products.CreateProductRequest{
		Name:        "Example Product",
		Description: "This is an example product",
		Price:       "99.99",
		CategoryId:  "category-1",
	}
	createResp, err := a.prodServiceClient.CreateProduct(ctx, createReq)
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Errorln("Failed to create product")
		return err
	}
	logging.GetLogger(ctx).Infof("Created product with ID: %s", createResp.Product.Id)
	
	logging.GetLogger(ctx).Infoln("Product service call completed")
	return nil
}

// testProductServiceConnection tests the connection to product service
func testProductServiceConnection(ctx context.Context, client pb_prod_products.ProductServiceClient) error {
	// Создаем контекст с таймаутом для тестирования соединения
	testCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	// Пытаемся получить список продуктов
	req := &pb_prod_products.AllProductsRequest{}
	_, err := client.AllProducts(testCtx, req)
	if err != nil {
		return fmt.Errorf("failed to test product service connection: %w", err)
	}
	
	logging.GetLogger(ctx).Infoln("Product service connection test successful")
	return nil
}


