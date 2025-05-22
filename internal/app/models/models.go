package models

import (
	"time"

	"github.com/google/uuid"
)

// Project represents a project with environment variables
type Project struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}

// Environment represents an environment type (development, sit, uat, etc.)
type Environment struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// EnvVariable represents a single environment variable
type EnvVariable struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	ProjectID     uuid.UUID  `db:"project_id" json:"project_id"`
	EnvironmentID uuid.UUID  `db:"environment_id" json:"environment_id"`
	Key           string     `db:"key" json:"key"`
	Value         string     `db:"value" json:"value"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at"`
}

// ProjectWithEnv represents a project with its environment variables
type ProjectWithEnv struct {
	Project      Project
	Environment  Environment
	EnvVariables []EnvVariable
}
