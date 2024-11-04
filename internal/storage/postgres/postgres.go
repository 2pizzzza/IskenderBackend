package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	"github.com/2pizzzza/plumbing/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(cfg *config.Config) (DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	dbConn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return DB{}, fmt.Errorf("unable to create connection pool: %v", err)
	}

	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		return DB{}, fmt.Errorf("unable to open sql database connection: %v", err)
	}
	defer sqlDB.Close()

	m, err := migrate.New(
		"file://database/migration",
		connStr)
	if err != nil {
		log.Printf("migration setup error: %s", err)
		return DB{}, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("migration error: %s", err)
		return DB{}, err
	}

	slog.Info("Success connection to database")
	return DB{Pool: dbConn}, nil
}
