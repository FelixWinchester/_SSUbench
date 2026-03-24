package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	App      AppConfig
}

type ServerConfig struct {
	Port         string        `env:"SERVER_PORT"          env-default:"8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"  env-default:"10s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT"  env-default:"60s"`
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"     env-default:"localhost"`
	Port     string `env:"DB_PORT"     env-default:"5432"`
	User     string `env:"DB_USER"     env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"postgres"`
	Name     string `env:"DB_NAME"     env-default:"ssubench"`
	SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
}

type JWTConfig struct {
	Secret string        `env:"JWT_SECRET" env-required:"true"`
	TTL    time.Duration `env:"JWT_TTL"    env-default:"24h"`
}

type AppConfig struct {
	LogLevel string `env:"LOG_LEVEL" env-default:"debug"`
}

func (d *DatabaseConfig) DSN() string {
	return "postgres://" + d.User + ":" + d.Password +
		"@" + d.Host + ":" + d.Port +
		"/" + d.Name + "?sslmode=" + d.SSLMode
}

func MustLoad() *Config {
	// Загружаем .env файл если он есть (в проде его не будет)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return cfg
}
