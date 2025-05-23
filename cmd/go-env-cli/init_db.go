package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go-env-cli/config"
	"go-env-cli/internal/pkg/db"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	fmt.Println("Loading configuration...")
	cfg, err := config.LoadConfig(".go-env-cli.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	dbConn, err := db.NewDB(db.Config{GO_CLI_DB: cfg.GO_CLI_DB})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Get migration files
	// Try to find migrations directory
	possiblePaths := []string{
		filepath.Join(".", "db", "migrations"),
		filepath.Join("..", "..", "db", "migrations"),
		filepath.Join(os.Getenv("HOME"), "go-env-cli", "db", "migrations"),
	}

	var migrationsDir string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			migrationsDir = path
			break
		}
	}

	if migrationsDir == "" {
		log.Fatalf("Could not find migrations directory in any of the expected locations")
	}

	fmt.Printf("Running migrations from %s...\n", migrationsDir)

	// Initialize migration manager
	migrationManager, err := db.NewMigrationManager(dbConn, migrationsDir)
	if err != nil {
		log.Fatalf("Failed to initialize migration manager: %v", err)
	}

	// Run migrations
	if err := migrationManager.MigrateUp(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Database initialization complete!")
}
