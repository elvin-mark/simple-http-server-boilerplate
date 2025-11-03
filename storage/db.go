package storage

import (
	"context"
	"fmt"
	"http-server/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

// InitDB initializes the database connection.
func InitDB(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := createUsersTable(pool); err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	return pool, nil
}

func createUsersTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE
		)
	`)
	return err
}
