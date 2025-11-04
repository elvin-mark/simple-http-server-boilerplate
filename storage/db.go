package storage

import (
	"context"
	"fmt"
	"http-server/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

// InitDB initializes the database connection.
func InitDB(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return pool, nil
}