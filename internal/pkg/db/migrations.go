package db

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MigrationManager struct {
	db         *sqlx.DB
	migrations []string
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sqlx.DB, migrationsFolder string) (*MigrationManager, error) {
	// Read migration files from the filesystem
	var migrations []string
	err := filepath.Walk(migrationsFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".sql") {
			migrations = append(migrations, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error reading migration files: %w", err)
	}

	// Sort migrations by filename
	sort.Strings(migrations)

	return &MigrationManager{
		db:         db,
		migrations: migrations,
	}, nil
}

// MigrateUp executes all migration files
func (m *MigrationManager) MigrateUp() error {
	// Create migrations table if it doesn't exist
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating migrations table: %w", err)
	}

	// Check which migrations have been applied
	appliedMigrations := make(map[string]bool)
	rows, err := m.db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("error querying applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("error scanning migration version: %w", err)
		}
		appliedMigrations[version] = true
	}

	// Apply each migration
	for _, migrationPath := range m.migrations {
		version := filepath.Base(migrationPath)
		if appliedMigrations[version] {
			log.Printf("Migration %s already applied, skipping", version)
			continue
		}

		log.Printf("Applying migration: %s", version)

		// Read migration content
		content, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %w", version, err)
		}

		// Execute migration in a transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("error starting transaction: %w", err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing migration %s: %w", version, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("error recording migration %s: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing migration %s: %w", version, err)
		}

		log.Printf("Successfully applied migration: %s", version)
	}

	return nil
}
