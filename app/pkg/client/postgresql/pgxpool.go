package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Client интерфейс для работы с PostgreSQL
type Client interface {
	Close()
	Ping(ctx context.Context) error
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// pgClient реализует интерфейс Client
type pgClient struct {
	pool *pgxpool.Pool
}

// Close закрывает соединение с базой данных
func (c *pgClient) Close() {
	c.pool.Close()
}

// Ping проверяет соединение с базой данных
func (c *pgClient) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

// Begin начинает транзакцию
func (c *pgClient) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.pool.Begin(ctx)
}

// BeginFunc выполняет функцию в транзакции
func (c *pgClient) BeginFunc(ctx context.Context, f func(pgx.Tx) error) error {
	return c.pool.BeginFunc(ctx, f)
}

// BeginTxFunc выполняет функцию в транзакции с опциями
func (c *pgClient) BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error {
	return c.pool.BeginTxFunc(ctx, txOptions, f)
}

// Query выполняет запрос и возвращает строки
func (c *pgClient) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, args...)
}

// QueryRow выполняет запрос и возвращает одну строку
func (c *pgClient) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.pool.QueryRow(ctx, sql, args...)
}

// Exec выполняет команду
func (c *pgClient) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, sql, arguments...)
}

type pgConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// NewPgConfig creates new pg config instance
func NewPgConfig(username string, password string, host string, port string, database string) *pgConfig {
	return &pgConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}
}

// NewClient
func NewClient(ctx context.Context, maxAttempts int, maxDelay time.Duration, cfg *pgConfig) (Client, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password,
		cfg.Host, cfg.Port, cfg.Database,
	)

	var pool *pgxpool.Pool
	err := DoWithAttempts(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pgxCfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			log.Fatalf("Unable to parse config: %v\n", err)
		}

		// pgxCfg.ConnConfig.Logger = logrusadapter.NewLogger(logger)

		pool, err = pgxpool.ConnectConfig(ctx, pgxCfg)
		if err != nil {
			log.Println("Failed to connect to postgres... Going to do the next attempt")

			return err
		}

		return nil
	}, maxAttempts, maxDelay)

	if err != nil {
		log.Fatal("All attempts are exceeded. Unable to connect to postgres")
	}

	return &pgClient{pool: pool}, nil
}

func DoWithAttempts(fn func() error, maxAttempts int, delay time.Duration) error {
	var err error

	for maxAttempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			maxAttempts--

			continue
		}

		return nil
	}

	return err
}

