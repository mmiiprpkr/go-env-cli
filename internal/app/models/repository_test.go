package models

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// getTestDB returns a connection to a test database
func getTestDB(t *testing.T) *sqlx.DB {
	// Get database connection parameters from environment variables
	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("TEST_DB_PORT")
	if port == "" {
		port = "5434"
	}
	user := os.Getenv("TEST_DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("TEST_DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("TEST_DB_NAME")
	if dbname == "" {
		dbname = "go-env-cli-test"
	}

	// Skip tests if CI environment variable is set but no database is available
	if os.Getenv("CI") != "" && (host == "" || port == "" || user == "" || password == "" || dbname == "") {
		t.Skip("Skipping database tests in CI environment with missing database configuration")
	}

	// Connect to database
	dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Skipf("Skipping database tests: %v", err)
		return nil
	}

	return db
}

func TestRepository_CreateProject(t *testing.T) {
	// Get test database connection
	db := getTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	// Create repository
	repo := NewRepository(db)

	// Test creating a project
	projectName := "test-project-" + uuid.New().String()
	project, err := repo.CreateProject(projectName, "Test project description")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Verify project was created
	if project.ID == uuid.Nil {
		t.Error("Expected project ID to be set")
	}
	if project.Name != projectName {
		t.Errorf("Expected project name to be %q, got %q", projectName, project.Name)
	}
	if project.Description != "Test project description" {
		t.Errorf("Expected project description to be %q, got %q", "Test project description", project.Description)
	}
	if project.CreatedAt.IsZero() {
		t.Error("Expected project creation date to be set")
	}
	if project.UpdatedAt.IsZero() {
		t.Error("Expected project update date to be set")
	}
	if project.DeletedAt != nil {
		t.Error("Expected project deletion date to be nil")
	}

	// Clean up
	_, err = db.Exec("DELETE FROM projects WHERE id = $1", project.ID)
	if err != nil {
		t.Logf("Warning: Failed to clean up test project: %v", err)
	}
}

func TestRepository_GetProjectByName(t *testing.T) {
	// Get test database connection
	db := getTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	// Create repository
	repo := NewRepository(db)

	// Create a project for testing
	projectName := "test-get-project-" + uuid.New().String()
	createdProject, err := repo.CreateProject(projectName, "Test get project")
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	// Test getting the project
	project, err := repo.GetProjectByName(projectName)
	if err != nil {
		t.Fatalf("Failed to get project by name: %v", err)
	}

	// Verify project fields
	if project.ID != createdProject.ID {
		t.Errorf("Expected project ID %v, got %v", createdProject.ID, project.ID)
	}
	if project.Name != projectName {
		t.Errorf("Expected project name %q, got %q", projectName, project.Name)
	}

	// Test getting a non-existent project
	_, err = repo.GetProjectByName("non-existent-project")
	if err == nil {
		t.Error("Expected error when getting non-existent project, got nil")
	}

	// Clean up
	_, err = db.Exec("DELETE FROM projects WHERE id = $1", project.ID)
	if err != nil {
		t.Logf("Warning: Failed to clean up test project: %v", err)
	}
}

func TestRepository_EnvVariableOperations(t *testing.T) {
	// Get test database connection
	db := getTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	// Create repository
	repo := NewRepository(db)

	// Create a project for testing
	projectName := "test-env-vars-" + uuid.New().String()
	project, err := repo.CreateProject(projectName, "Test env variables")
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	// Get development environment
	env, err := repo.GetEnvironmentByName("development")
	if err != nil {
		t.Fatalf("Failed to get development environment: %v", err)
	}

	// Create a key-value pair
	key := "TEST_KEY"
	value := "test_value"
	variable, err := repo.SetEnvVariable(project.ID, env.ID, key, value)
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	// Verify variable was set
	if variable.Key != key {
		t.Errorf("Expected variable key %q, got %q", key, variable.Key)
	}
	if variable.Value != value {
		t.Errorf("Expected variable value %q, got %q", value, variable.Value)
	}

	// Get the variable
	getVar, err := repo.GetEnvVariable(project.ID, env.ID, key)
	if err != nil {
		t.Fatalf("Failed to get environment variable: %v", err)
	}
	if getVar.Key != key || getVar.Value != value {
		t.Errorf("Got unexpected variable: key=%q value=%q", getVar.Key, getVar.Value)
	}

	// Update the variable
	newValue := "updated_value"
	updatedVar, err := repo.SetEnvVariable(project.ID, env.ID, key, newValue)
	if err != nil {
		t.Fatalf("Failed to update environment variable: %v", err)
	}
	if updatedVar.Value != newValue {
		t.Errorf("Expected updated value %q, got %q", newValue, updatedVar.Value)
	}

	// Delete the variable
	err = repo.DeleteEnvVariable(project.ID, env.ID, key)
	if err != nil {
		t.Fatalf("Failed to delete environment variable: %v", err)
	}

	// Verify it was deleted (soft delete)
	_, err = repo.GetEnvVariable(project.ID, env.ID, key)
	if err == nil {
		t.Error("Expected error when getting deleted variable, got nil")
	}

	// Clean up
	_, err = db.Exec("DELETE FROM env_variables WHERE project_id = $1", project.ID)
	if err != nil {
		t.Logf("Warning: Failed to clean up test env variables: %v", err)
	}
	_, err = db.Exec("DELETE FROM projects WHERE id = $1", project.ID)
	if err != nil {
		t.Logf("Warning: Failed to clean up test project: %v", err)
	}
}
