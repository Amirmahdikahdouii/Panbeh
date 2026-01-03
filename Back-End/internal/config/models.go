package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type LogLevel string

func NewLogLevel(level string) LogLevel {
	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		return LogLevelDebug
	case "INFO":
		return LogLevelInfo
	case "WARN":
		return LogLevelWarn
	case "ERROR":
		return LogLevelError
	}
	panic(fmt.Sprintf("Invalid log level in config: %s", level))
}

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

type PostgresDSN string

func NewPostgresDSN(dsn string) (PostgresDSN, error) {
	validateURLDSN := func(dsn string) (PostgresDSN, error) {
		u, err := url.Parse(dsn)
		if err != nil {
			return "", err
		}

		if u.Scheme != "postgres" && u.Scheme != "postgresql" {
			return "", errors.New("invalid scheme, expected postgres or postgresql")
		}

		if u.Host == "" {
			return "", errors.New("host is missing")
		}

		host, _, err := net.SplitHostPort(u.Host)
		if err != nil {
			// allow host without port
			host = u.Host
		}

		if host == "" {
			return "", errors.New("host is empty")
		}

		if strings.Trim(u.Path, "/") == "" {
			return "", errors.New("database name is missing")
		}

		return PostgresDSN(dsn), nil
	}

	validateKeyValueDSN := func(dsn string) (PostgresDSN, error) {
		parts := strings.Fields(dsn)
		if len(parts) == 0 {
			return "", errors.New("invalid DSN format")
		}

		kv := make(map[string]string)
		for _, part := range parts {
			split := strings.SplitN(part, "=", 2)
			if len(split) != 2 {
				return "", errors.New("invalid key=value pair: " + part)
			}
			kv[split[0]] = split[1]
		}

		// Minimal required fields
		if kv["host"] == "" {
			return "", errors.New("host is required")
		}
		if kv["dbname"] == "" {
			return "", errors.New("dbname is required")
		}

		return PostgresDSN(dsn), nil
	}

	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return "", errors.New("dsn is empty")
	}

	// Case 1: URL-style DSN
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		return validateURLDSN(dsn)
	}

	// Case 2: keyword/value DSN
	return validateKeyValueDSN(dsn)
}

func NewPostgresConfig(postgresDSN string, maxOpenConns int, maxIdleConns int, connMaxLifetime time.Duration) Postgres {
	dsn, err := NewPostgresDSN(postgresDSN)
	if err != nil {
		panic(err)
	}

	return Postgres{
		PostgresDSN:             dsn,
		PostgresMaxOpenConns:    maxOpenConns,
		PostgresMaxIdleConns:    maxIdleConns,
		PostgresConnMaxLifetime: connMaxLifetime,
	}
}

type RedisAddr string

func NewRedisAddr(addr string) (RedisAddr, error) {
	validateRedisURL := func(addr string) (RedisAddr, error) {
		u, err := url.Parse(addr)
		if err != nil {
			return "", err
		}

		if u.Scheme != "redis" && u.Scheme != "rediss" {
			return "", errors.New("invalid scheme, expected redis or rediss")
		}

		if u.Host == "" {
			return "", errors.New("host is missing")
		}

		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return "", errors.New("invalid host:port")
		}

		if host == "" {
			return "", errors.New("host is empty")
		}

		if port == "" {
			return "", errors.New("port is empty")
		}

		if _, err := strconv.Atoi(port); err != nil {
			return "", errors.New("port must be a number")
		}

		// Optional DB index validation (/0, /1, ...)
		if u.Path != "" && u.Path != "/" {
			db := strings.TrimPrefix(u.Path, "/")
			if _, err := strconv.Atoi(db); err != nil {
				return "", errors.New("invalid redis db index")
			}
		}

		return RedisAddr(addr), nil
	}

	validateRedisHostPort := func(addr string) (RedisAddr, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return "", errors.New("invalid redis address, expected host:port")
		}

		if host == "" {
			return "", errors.New("host is empty")
		}

		if port == "" {
			return "", errors.New("port is empty")
		}

		if _, err := strconv.Atoi(port); err != nil {
			return "", errors.New("port must be a number")
		}

		return RedisAddr(addr), nil
	}
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "", errors.New("redis address is empty")
	}

	// URL-style DSN
	if strings.HasPrefix(addr, "redis://") || strings.HasPrefix(addr, "rediss://") {
		return validateRedisURL(addr)
	}

	// host:port style
	return validateRedisHostPort(addr)
}

func NewRedisConfig(addr string) Redis {
	redisAddr, err := NewRedisAddr(addr)
	if err != nil {
		panic(err)
	}

	return Redis{
		RedisAddr: redisAddr,
	}
}

type Config struct {
	Env      string
	LogLevel LogLevel

	// TODO: add a type for valid rest addr
	RestAddr string
	Postgres
	Redis
	OTPTTL time.Duration
}

type Postgres struct {
	PostgresDSN             PostgresDSN
	PostgresMaxOpenConns    int
	PostgresMaxIdleConns    int
	PostgresConnMaxLifetime time.Duration
}

type Redis struct {
	RedisAddr RedisAddr
}
