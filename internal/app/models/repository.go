package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for environment variables
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// CreateProject creates a new project
func (r *Repository) CreateProject(name, description string) (*Project, error) {
	project := &Project{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `
		INSERT INTO projects (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, created_at, updated_at
	`

	err := r.db.QueryRowx(query,
		project.ID,
		project.Name,
		project.Description,
		project.CreatedAt,
		project.UpdatedAt,
	).StructScan(project)

	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

// GetProjectByName retrieves a project by name
func (r *Repository) GetProjectByName(name string) (*Project, error) {
	project := &Project{}
	query := `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM projects
		WHERE name = $1 AND deleted_at IS NULL
	`

	err := r.db.Get(project, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get project by name: %w", err)
	}

	return project, nil
}

// GetAllProjects retrieves all non-deleted projects
func (r *Repository) GetAllProjects() ([]Project, error) {
	projects := []Project{}
	query := `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM projects
		WHERE deleted_at IS NULL
		ORDER BY name
	`

	err := r.db.Select(&projects, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all projects: %w", err)
	}

	return projects, nil
}

// SearchProjects searches for projects by name pattern
func (r *Repository) SearchProjects(pattern string) ([]Project, error) {
	projects := []Project{}
	query := `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM projects
		WHERE name ILIKE $1 AND deleted_at IS NULL
		ORDER BY name
	`

	err := r.db.Select(&projects, query, "%"+pattern+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search projects: %w", err)
	}

	return projects, nil
}

// SoftDeleteProject soft-deletes a project
func (r *Repository) SoftDeleteProject(id uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE projects
		SET deleted_at = $1, updated_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, now, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no project found with ID %s", id)
	}

	return nil
}

// GetEnvironmentByName retrieves an environment by name
func (r *Repository) GetEnvironmentByName(name string) (*Environment, error) {
	env := &Environment{}
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM environments
		WHERE name = $1
	`

	err := r.db.Get(env, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment by name: %w", err)
	}

	return env, nil
}

// GetAllEnvironments retrieves all environments
func (r *Repository) GetAllEnvironments() ([]Environment, error) {
	environments := []Environment{}
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM environments
		ORDER BY name
	`

	err := r.db.Select(&environments, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all environments: %w", err)
	}

	return environments, nil
}

// CreateEnvironment creates a new environment
func (r *Repository) CreateEnvironment(name, description string) (*Environment, error) {
	env := &Environment{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `
		INSERT INTO environments (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (name) DO NOTHING
		RETURNING id, name, description, created_at, updated_at
	`

	err := r.db.QueryRowx(query,
		env.ID,
		env.Name,
		env.Description,
		env.CreatedAt,
		env.UpdatedAt,
	).StructScan(env)

	if err != nil {
		return nil, fmt.Errorf("failed to create environment: %w", err)
	}

	return env, nil
}

// SetEnvVariable sets (creates or updates) an environment variable
func (r *Repository) SetEnvVariable(projectID, environmentID uuid.UUID, key, value string) (*EnvVariable, error) {
	now := time.Now()
	query := `
		INSERT INTO env_variables (id, project_id, environment_id, key, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (project_id, environment_id, key) DO UPDATE
		SET value = $5, updated_at = $7
		RETURNING id, project_id, environment_id, key, value, created_at, updated_at, deleted_at
	`

	variable := &EnvVariable{
		ID:            uuid.New(),
		ProjectID:     projectID,
		EnvironmentID: environmentID,
		Key:           key,
		Value:         value,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err := r.db.QueryRowx(query,
		variable.ID,
		variable.ProjectID,
		variable.EnvironmentID,
		variable.Key,
		variable.Value,
		variable.CreatedAt,
		variable.UpdatedAt,
	).StructScan(variable)

	if err != nil {
		return nil, fmt.Errorf("failed to set environment variable: %w", err)
	}

	return variable, nil
}

// GetEnvVariable gets an environment variable by key
func (r *Repository) GetEnvVariable(projectID, environmentID uuid.UUID, key string) (*EnvVariable, error) {
	variable := &EnvVariable{}
	query := `
		SELECT id, project_id, environment_id, key, value, created_at, updated_at, deleted_at
		FROM env_variables
		WHERE project_id = $1 AND environment_id = $2 AND key = $3 AND deleted_at IS NULL
	`

	err := r.db.Get(variable, query, projectID, environmentID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment variable: %w", err)
	}

	return variable, nil
}

// GetEnvVariables gets all environment variables for a project and environment
func (r *Repository) GetEnvVariables(projectID, environmentID uuid.UUID) ([]EnvVariable, error) {
	variables := []EnvVariable{}
	query := `
		SELECT id, project_id, environment_id, key, value, created_at, updated_at, deleted_at
		FROM env_variables
		WHERE project_id = $1 AND environment_id = $2 AND deleted_at IS NULL
		ORDER BY key
	`

	err := r.db.Select(&variables, query, projectID, environmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment variables: %w", err)
	}

	return variables, nil
}

// DeleteEnvVariable deletes an environment variable
func (r *Repository) DeleteEnvVariable(projectID, environmentID uuid.UUID, key string) error {
	now := time.Now()
	query := `
		UPDATE env_variables
		SET deleted_at = $1, updated_at = $1
		WHERE project_id = $2 AND environment_id = $3 AND key = $4 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, now, projectID, environmentID, key)
	if err != nil {
		return fmt.Errorf("failed to delete environment variable: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no environment variable found with key %s", key)
	}

	return nil
}

// GetEnvironmentsForProject retrieves all environments used by a specific project
func (r *Repository) GetEnvironmentsForProject(projectID uuid.UUID) ([]Environment, error) {
	environments := []Environment{}
	query := `
		SELECT DISTINCT e.id, e.name, e.description, e.created_at, e.updated_at
		FROM environments e
		JOIN env_variables ev ON e.id = ev.environment_id
		WHERE ev.project_id = $1 AND ev.deleted_at IS NULL
		ORDER BY e.name
	`

	err := r.db.Select(&environments, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get environments for project: %w", err)
	}

	return environments, nil
}
