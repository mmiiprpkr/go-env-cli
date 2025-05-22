# go-env-cli

A powerful CLI tool for managing environment variables across multiple projects and environments, with PostgreSQL storage.

## Features

- **Import .env files** - Import variables from .env files into the database
- **Export .env files** - Export stored variables to .env files
- **Project Management** - List, search, and soft-delete projects
- **Environment Variables** - Set, get, delete, and list environment variables
- **Environment Grouping** - Support for multiple environments (development, sit, uat, etc.)

## Installation

### Prerequisites

- Go 1.18 or higher
- PostgreSQL database

### Building from Source

Clone the repository and build the project:

```bash
git clone https://github.com/yourusername/go-env-cli.git
cd go-env-cli
go build -o go-env-cli .
```

## Configuration

Configuration can be provided via:

1. Config file (`.go-env-cli.yaml` in the current directory or home)
2. Environment variables (prefixed with `GO_ENV_CLI_`)
3. Command line flags

### Database Connection

Set up the database connection via environment variables:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_env_cli
```

## Initializing the Database

Before first use, initialize the database:

```bash
go run cmd/go-env-cli/init_db.go
```

## Usage

### Basic Commands

```bash
# Import variables from a .env file
go-env-cli import .env --project my-project --env development

# Export variables to a .env file
go-env-cli export .env.production --project my-project --env production

# List all projects (now includes environment information)
go-env-cli list-projects

# Get detailed project information including environments
go-env-cli project-details --project my-project

# Search for projects
go-env-cli search-project api

# Set an environment variable
go-env-cli set --project my-project --env development --key API_KEY --value "secret123"

# Get an environment variable
go-env-cli get --project my-project --env development --key API_KEY

# Delete an environment variable
go-env-cli delete --project my-project --env development --key API_KEY

# List all environment variables for a project
go-env-cli list --project my-project --env development

# Soft delete a project
go-env-cli delete-project --project old-project

# List all environments
go-env-cli env list

# Create a new environment
go-env-cli env create --name staging --description "Staging environment"
```

### Using with Docker

You can use the CLI with a Docker-based PostgreSQL for portability:

```bash
docker run -d --name postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=go_env_cli -p 5432:5432 postgres:14
```

## License

MIT
