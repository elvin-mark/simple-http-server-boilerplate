package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"http-server/config"

	_ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL driver
)

// Run applies all pending database migrations.
func Run(cfg *config.DatabaseConfig) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := ensureSchemaMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}

	appliedVersions, err := getAppliedVersions(db)
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	files, err := filepath.Glob("migrations/*.up.sql")
	if err != nil {
		return fmt.Errorf("failed to find migration files: %w", err)
	}
	sort.Strings(files)

	for _, file := range files {
		version, err := getVersionFromFile(file)
		if err != nil {
			log.Printf("Skipping file with invalid version format: %s", file)
			continue
		}

		if !appliedVersions[version] {
			log.Printf("Applying migration %s...", file)
			if err := applyMigration(db, file, version); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", file, err)
			}
			log.Printf("Successfully applied migration %s", file)
		}
	}

	return nil
}

func ensureSchemaMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY
		)
	`)
	return err
}

func getAppliedVersions(db *sql.DB) (map[int64]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make(map[int64]bool)
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions[version] = true
	}
	return versions, nil
}

func getVersionFromFile(file string) (int64, error) {
	base := filepath.Base(file)
	parts := strings.Split(base, "_")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid migration filename format")
	}
	return strconv.ParseInt(parts[0], 10, 64)
}

func applyMigration(db *sql.DB, file string, version int64) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		// Rollback on error
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err := tx.Exec(string(content)); err != nil {
		return err
	}

	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return err
	}

	return tx.Commit()
}
