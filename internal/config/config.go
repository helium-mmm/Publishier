package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Server ServerConfig
	DB DBConfig
	Telegram TelegramConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	DSN string
	MaxOpenConns int
	MaxIdleConns int
}

type TelegramConfig struct {
	BaseURL string
	Timeout time.Duration
}

func MustLoad() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		DB: DBConfig{
			DSN: mustEnv("DB_DSN"),
			MaxOpenConns: 20,
			MaxIdleConns: 10,
		},
		Telegram: TelegramConfig{
			BaseURL: "https://api.telegram.org",
			Timeout: 10 * time.Second,
		},
	}

	return cfg
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal("missing env: ", key)
	}
	return val
}

func getEnv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}