package databases

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/panbeh/otp-backend/internal/config"
)

func ConnectPostgres(cfg *config.Postgres) *sql.DB {
	db, err := sql.Open("pgx", string(cfg.PostgresDSN))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.PostgresMaxOpenConns)
	db.SetMaxIdleConns(cfg.PostgresMaxIdleConns)
	db.SetConnMaxLifetime(cfg.PostgresConnMaxLifetime)

	if err := db.PingContext(context.Background()); err != nil {
		panic(err)
	}
	return db
}
