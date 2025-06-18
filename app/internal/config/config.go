package config

import (
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config содержит конфигурацию приложения
type Config struct {
	IsDebug bool  `yaml:"is-debug" env:"IS_DEBUG" env-default:"false"` 
	IsDevelopment bool `yaml:"is-development" env:"IS_DEV" env-default:"false"`
	HTTP struct {
		IP string `yaml:"ip" env:"HTTP-IP"`
		Port int `yaml:"port" env:"HTTP-PORT"`
		ReadTimeout time.Duration `yaml:"read-timeout" env:"HTTP-READ_TIMEOUT"`
		WriteTimeout time.Duration `yaml:"write-timeout" env:"HTTP-WRITE_TIMEOUT"`
		CORS struct {
			AllowedMethods     []string `yaml:"allowed_methods" env:"HTTP-CORS-ALLOWED-METHODS"`
			AllowedOrigins     []string `yaml:"allowed_origins"`
			AllowCredentials   bool     `yaml:"allow_credentials"`
			AllowedHeaders     []string `yaml:"allowed_headers"`
			OptionsPassthrough bool     `yaml:"options_passthrough"`
			ExposedHeaders     []string `yaml:"exposed_headers"`
			Debug              bool     `yaml:"debug"`
		} `yaml:"cors"`
	} `yaml:"http"`
	GRPC struct {
		IP   string `yaml:"ip" env:"GRPC-IP"`
		Port int    `yaml:"port" env:"GRPC-PORT"`
	} `yaml:"grpc"`
	AppConfig struct {
		LogLevel string `yaml:"log-level" env:"LOG_LEVEL" env-default:"trace" env-description:"Log level for the application. Options: debug, info, warn, error, fatal"`
		AdminUser struct {
			Email    string `yaml:"email" env:"ADMIN_EMAIL" env-default:"admin"`
			Password string `yaml:"password" env:"ADMIN_PASSWORD" env-default:"admin"`
		} `yaml:"admin"`
	} `yaml:"app"`
	PostgreSQL struct {
		Username string `yaml:"username" env:"PSQL_USERNAME" env-required:"true"`
		Password string `yaml:"password" env:"PSQL_PASSWORD" env-required:"true"`
		Host     string `yaml:"host" env:"PSQL_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"PSQL_PORT" env-required:"true"`
		Database string `yaml:"database" env:"PSQL_DATABASE" env-required:"true"`
	} `yaml:"postgresql"`
}

const (
	EnvConfigPathName = "CONFIG_PATH"
	FlagConfigPathName = "config"
)

var configPath string
var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		flag.StringVar(&configPath, FlagConfigPathName, "configs/config.local.yaml", "this is app config file")
		flag.Parse()

		log.Print("config init")

		if configPath == "" {
			configPath = os.Getenv(EnvConfigPathName)
		}

		if configPath == "" {
			log.Fatal("config is required")
		}

		instance = &Config{}

		// Сначала читаем YAML файл
		if err := cleanenv.ReadConfig(configPath, instance); err != nil {
			log.Printf("Error reading config file %s: %v", configPath, err)
		}

		// Затем читаем .env файл (если есть)
		if err := cleanenv.ReadConfig(".env", instance); err != nil {
			// Если .env файл не найден, пробуем только переменные окружения
			if err := cleanenv.ReadEnv(instance); err != nil {
				helpText := "The Art of Development - Production Service"
				help, _ := cleanenv.GetDescription(instance, &helpText)
				log.Print(help)
				log.Fatal(err)
			}
		} 
	})

	return instance
}


