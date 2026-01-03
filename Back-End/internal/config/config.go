package config

import (
	"os"
	"strconv"
	"time"
)

var cfg *Config

func Load() error {
	cfg = &Config{
		Env:      getenv("APP_ENV", "local"),
		LogLevel: NewLogLevel(getenv("LOG_LEVEL", "info")),
		HTTPPort: getenv("HTTP_PORT", "8080"),
	}

	cfg.Postgres = NewPostgresConfig(getenv("POSTGRES_DSN", "postgres://admin:123456@localhost:5432/panbeh"), parseInt(getenv("POSTGRES_MAX_OPEN_CONNS", "10")), parseInt(getenv("POSTGRES_MAX_IDLE_CONNS", "5")), parseIntDuration(getenv("POSTGRES_CONN_MAX_LIFETIME_SECONDS", "300")))
	cfg.Redis = NewRedisConfig(getenv("REDIS_ADDR", "127.0.0.1:6379"))
	cfg.OTPTTL = parseIntDuration(getenv("OTP_TTL_SECONDS", "300"))
	return nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseInt(v string) int {
	num, err := strconv.Atoi(v)
	if err != nil {
		panic("Failed to parse int: " + v)
	}
	return num
}

func parseIntDuration(v string) time.Duration {
	num, err := strconv.Atoi(v)
	if err != nil {
		panic("Failed to parse int: " + v)
	}
	return time.Duration(num) * time.Second
}

func GetLogLevel() LogLevel {
	return cfg.LogLevel
}

func GetPostgres() *Postgres {
	return &cfg.Postgres
}
