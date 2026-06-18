package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Telegram TelegramConfig
	JWT      JWTConfig
	Crypto   CryptoConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
}

type TelegramConfig struct {
	BaseURL string
	Timeout time.Duration
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

type CryptoConfig struct {
	Key string
}

func MustLoad() *Config {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = buildDSN()
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		DB: DBConfig{
			DSN:          dsn,
			MaxOpenConns: 20,
			MaxIdleConns: 10,
		},
		Telegram: TelegramConfig{
			BaseURL: getEnv("TELEGRAM_BASE_URL", "https://api.telegram.org"),
			Timeout: 10 * time.Second,
		},
		JWT: JWTConfig{
			Secret: mustEnv("JWT_SECRET"),
			TTL:    24 * time.Hour,
		},
		Crypto: CryptoConfig{
			Key: mustEnv("AES_SECRET"),
		},
	}
}

func buildDSN() string {
	host := mustEnv("DB_HOST")
	port := getEnv("DB_PORT", "5432")
	user := mustEnv("DB_USER")
	password := mustEnv("DB_PASSWORD")
	dbName := getEnv("DB_NAME", "publishier")
	sslMode := getEnv("DB_SSLMODE", "disable")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, sslMode,
	)
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
