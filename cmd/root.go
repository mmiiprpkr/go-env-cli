package cmd

import (
	"fmt"
	"os"
	"strconv"

	"go-env-cli/config"
	"go-env-cli/internal/app/handlers"
	"go-env-cli/internal/app/models"
	"go-env-cli/internal/pkg/db"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	configFile string
	envFile    string

	// Command specific flags
	projectName     string
	environmentName string
	keyName         string
	keyValue        string
	description     string
	force           bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-env-cli",
	Short: "A CLI tool for managing environment variables",
	Long: `go-env-cli is a command-line tool that helps you manage environment variables
across multiple projects and environments. It stores variables in a PostgreSQL database
and provides commands for importing/exporting .env files, setting/getting variables,
and more.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load .env file if specified
		if envFile != "" {
			err := godotenv.Load(envFile)
			if err != nil {
				fmt.Printf("Error loading .env file: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ./.go-env-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&envFile, "env-file", "", "env file to load configuration from")

	// Add commands
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(listProjectsCmd)
	rootCmd.AddCommand(searchProjectCmd)
	rootCmd.AddCommand(setEnvCmd)
	rootCmd.AddCommand(getEnvCmd)
	rootCmd.AddCommand(deleteEnvCmd)
	rootCmd.AddCommand(listEnvCmd)
	rootCmd.AddCommand(softDeleteProjectCmd)
	rootCmd.AddCommand(environmentCmd)
	rootCmd.AddCommand(projectDetailsCmd)
}

// initHandler creates and initializes the environment handler
func initHandler() (*handlers.EnvHandler, error) {
	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	// If env file was specified, override with its values
	if envFile != "" {
		env, err := godotenv.Read(envFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read .env file: %v", err)
		}

		if host, ok := env["DB_HOST"]; ok {
			cfg.Database.Host = host
		}
		if port, ok := env["DB_PORT"]; ok {
			portNum, err := strconv.Atoi(port)
			if err == nil {
				cfg.Database.Port = portNum
			}
		}
		if user, ok := env["DB_USER"]; ok {
			cfg.Database.User = user
		}
		if password, ok := env["DB_PASSWORD"]; ok {
			cfg.Database.Password = password
		}
		if dbName, ok := env["DB_NAME"]; ok {
			cfg.Database.DBName = dbName
		}
	}

	// Connect to database
	dbConn, err := db.NewDB(db.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Create repository
	repo := models.NewRepository(dbConn)

	// Create handler
	handler := handlers.NewEnvHandler(repo)

	return handler, nil
}

// Import command
var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import environment variables from a .env file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		// Validate flags
		if projectName == "" {
			fmt.Println("Error: --project flag is required")
			os.Exit(1)
		}
		if environmentName == "" {
			environmentName = "development" // Default to development
		}

		// Initialize handler
		handler, err := initHandler()
		if err != nil {
			fmt.Printf("Error initializing: %v\n", err)
			os.Exit(1)
		}

		// Import file
		err = handler.ImportEnvFile(filePath, projectName, environmentName)
		if err != nil {
			fmt.Printf("Error importing .env file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully imported environment variables from %s to project '%s' (%s environment)\n",
			filePath, projectName, environmentName)
	},
}

// Export command
var exportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Export environment variables to a .env file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		// Validate flags
		if projectName == "" {
			fmt.Println("Error: --project flag is required")
			os.Exit(1)
		}
		if environmentName == "" {
			environmentName = "development" // Default to development
		}

		// Check if file exists and confirm overwrite if needed
		if _, err := os.Stat(filePath); err == nil {
			if !force && !cmd.Flags().Changed("force") {
				fmt.Printf("File %s already exists. Overwrite? [y/N]: ", filePath)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Export cancelled")
					return
				}
			}
		}

		// Initialize handler
		handler, err := initHandler()
		if err != nil {
			fmt.Printf("Error initializing: %v\n", err)
			os.Exit(1)
		}

		// Export to file
		err = handler.ExportEnvFile(filePath, projectName, environmentName)
		if err != nil {
			fmt.Printf("Error exporting to .env file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully exported environment variables from project '%s' (%s environment) to %s\n",
			projectName, environmentName, filePath)
	},
}

// List projects command
var listProjectsCmd = &cobra.Command{
	Use:   "list-projects",
	Short: "List all projects",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize handler
		handler, err := initHandler()
		if err != nil {
			fmt.Printf("Error initializing: %v\n", err)
			os.Exit(1)
		}

		// Get projects
		projects, err := handler.ListProjects()
		if err != nil {
			fmt.Printf("Error listing projects: %v\n", err)
			os.Exit(1)
		}

		// Display projects
		if len(projects) == 0 {
			fmt.Println("No projects found")
			return
		}

		fmt.Println("Projects:")
		fmt.Println("=========")
		for _, p := range projects {
			fmt.Printf("- %s: %s\n", p.Name, p.Description)

			// Get environments for this project
			environments, err := handler.GetEnvironmentsForProject(p.Name)
			if err == nil && len(environments) > 0 {
				fmt.Printf("  Environments: ")
				for i, env := range environments {
					if i > 0 {
						fmt.Printf(", ")
					}
					fmt.Printf("%s", env.Name)
				}
				fmt.Println()
			}
		}
	},
}

// Search project command
var searchProjectCmd = &cobra.Command{
	Use:   "search-project [pattern]",
	Short: "Search for projects by name pattern",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pattern := args[0]

		// Initialize handler
		handler, err := initHandler()
		if err != nil {
			fmt.Printf("Error initializing: %v\n", err)
			os.Exit(1)
		}

		// Search projects
		projects, err := handler.SearchProjects(pattern)
		if err != nil {
			fmt.Printf("Error searching projects: %v\n", err)
			os.Exit(1)
		}

		// Display projects
		if len(projects) == 0 {
			fmt.Printf("No projects found matching '%s'\n", pattern)
			return
		}

		fmt.Printf("Projects matching '%s':\n", pattern)
		fmt.Println("======================")
		for _, p := range projects {
			fmt.Printf("- %s: %s\n", p.Name, p.Description)
		}
	},
}

// Set env variable command
var setEnvCmd = &cobra.Command{
	Use:   "set",
	Short: "Set an environment variable",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate flags
		if projectName == "" {
			fmt.Println("Error: --project flag is required")
			os.Exit(1)
		}
		if environmentName == "" {
			environmentName = "development" // Default to development
		}
		if keyName == "" {
			fmt.Println("Error: --key flag is required")
			os.Exit(1)
		}

		// Initialize handler
		handler, err := initHandler()
		if err != nil {
			fmt.Printf("Error initializing: %v\n", err)
			os.Exit(1)
		}

		// Set variable
		err = handler.SetEnvVariable(projectName, environmentName, keyName, keyValue)
		if err != nil {
			fmt.Printf("Error setting environment variable: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully set %s=%s for project '%s' (%s environment)\n",
			keyName, keyValue, projectName, environmentName)
	},
}

// Get env variable command
var getEnvCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an environment variable",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate flags
		if projectName == "" {
			fmt.Println("Error: --project flag is required")
			os.Exit(1)
		}
		if environmentName == "" {
			environmentName = "development" // Default to development
		}
		if keyName == "" {
			fmt.Println("Error: --key flag is required")
			os.Exit(1)
		}

		// Initialize handler
		handler, err := initHandler()
		if err != nil {
			fmt.Printf("Error initializing: %v\n", err)
			os.Exit(1)
		}

		// Get variable
		value, err := handler.GetEnvVariable(projectName, environmentName, keyName)
		if err != nil {
			fmt.Printf("Error getting environment variable: %v\n", err)
			os.Exit(1)
		}

		// Just print the value (for piping to other commands)
		fmt.Println(value)
	},
}

// Delete env variable command
var deleteEnvCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an environment variable",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate flags
		if projectName == "" {
			fmt.Println("Error: --project flag is required")
			os.Exit(1)
		}
		if environmentName == "" {
			environmentName = "development" // Default to development
		}
		if keyName == "" {
			fmt.Println("Error: --key flag is required")
			os.Exit(1)
		}

		// Initialize handler
		handler, err := initHandler()
		if err != nil {
			fmt.Printf("Error initializing: %v\n", err)
			os.Exit(1)
		}

		// Delete variable
		err = handler.DeleteEnvVariable(projectName, environmentName, keyName)
		if err != nil {
			fmt.Printf("Error deleting environment variable: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully deleted environment variable '%s' from project '%s' (%s environment)\n",
			keyName, projectName, environmentName)
	},
}

// List env variables command
var listEnvCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environment variables for a project",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate flags
		if projectName == "" {
			fmt.Println("Error: --project flag is required")
			os.Exit(1)
		}
		if environmentName == "" {
			environmentName = "development" // Default to development
		}

		// Initialize handler
		handler, err := initHandler()
		if err != nil {
			fmt.Printf("Error initializing: %v\n", err)
			os.Exit(1)
		}

		// Get variables
		var variables []models.EnvVariable
		if keyName != "" {
			// Search by pattern
			variables, err = handler.SearchEnvVariables(projectName, environmentName, keyName)
		} else {
			// List all
			variables, err = handler.ListEnvVariables(projectName, environmentName)
		}

		if err != nil {
			fmt.Printf("Error listing environment variables: %v\n", err)
			os.Exit(1)
		}

		// Display variables
		if len(variables) == 0 {
			fmt.Printf("No environment variables found for project '%s' (%s environment)\n",
				projectName, environmentName)
			return
		}

		fmt.Printf("Environment variables for project '%s' (%s environment):\n",
			projectName, environmentName)
		fmt.Println("=================================================")
		for _, v := range variables {
			fmt.Printf("%s=%s\n", v.Key, v.Value)
		}
	},
}

// Soft delete project command
var softDeleteProjectCmd = &cobra.Command{
	Use:   "delete-project",
	Short: "Soft delete a project",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate flags
		if projectName == "" {
			fmt.Println("Error: --project flag is required")
			os.Exit(1)
		}

		// Confirm deletion unless --force is specified
		if !force && !cmd.Flags().Changed("force") {
			fmt.Printf("Are you sure you want to delete the project '%s'? This can't be undone. [y/N]: ", projectName)
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Delete cancelled")
				return
			}
		}

		// Initialize handler
		handler, err := initHandler()
		if err != nil {
			fmt.Printf("Error initializing: %v\n", err)
			os.Exit(1)
		}

		// Delete project
		err = handler.SoftDeleteProject(projectName)
		if err != nil {
			fmt.Printf("Error deleting project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully deleted project '%s'\n", projectName)
	},
}

// Environment command (with subcommands)
var environmentCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environments",
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior is to list environments
		cmd.Help()
	},
}

// List environments command
var listEnvironmentsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environments",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize handler
		handler, err := initHandler()
		if err != nil {
			fmt.Printf("Error initializing: %v\n", err)
			os.Exit(1)
		}

		// Get environments
		environments, err := handler.ListEnvironments()
		if err != nil {
			fmt.Printf("Error listing environments: %v\n", err)
			os.Exit(1)
		}

		// Display environments
		if len(environments) == 0 {
			fmt.Println("No environments found")
			return
		}

		fmt.Println("Environments:")
		fmt.Println("============")
		for _, e := range environments {
			fmt.Printf("- %s: %s\n", e.Name, e.Description)
		}
	},
}

// Create environment command
var createEnvironmentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate flags
		if environmentName == "" {
			fmt.Println("Error: --name flag is required")
			os.Exit(1)
		}

		// Initialize handler
		handler, err := initHandler()
		if err != nil {
			fmt.Printf("Error initializing: %v\n", err)
			os.Exit(1)
		}

		// Create environment
		err = handler.CreateEnvironment(environmentName, description)
		if err != nil {
			fmt.Printf("Error creating environment: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created environment '%s'\n", environmentName)
	},
}

// Show project details command
var projectDetailsCmd = &cobra.Command{
	Use:   "project-details",
	Short: "Show details of a project including its environments",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate flags
		if projectName == "" {
			fmt.Println("Error: --project flag is required")
			os.Exit(1)
		}

		// Initialize handler
		handler, err := initHandler()
		if err != nil {
			fmt.Printf("Error initializing: %v\n", err)
			os.Exit(1)
		}

		// Get projects to find the specific one
		projects, err := handler.ListProjects()
		if err != nil {
			fmt.Printf("Error listing projects: %v\n", err)
			os.Exit(1)
		}

		// Find the requested project
		var foundProject models.Project
		projectFound := false
		for _, p := range projects {
			if p.Name == projectName {
				foundProject = p
				projectFound = true
				break
			}
		}

		if !projectFound {
			fmt.Printf("Error: project '%s' not found\n", projectName)
			os.Exit(1)
		}

		// Get environments for the project
		environments, err := handler.GetEnvironmentsForProject(projectName)
		if err != nil {
			fmt.Printf("Error getting environments for project: %v\n", err)
			os.Exit(1)
		}

		// Display project details
		fmt.Printf("Project: %s\n", foundProject.Name)
		fmt.Printf("Description: %s\n", foundProject.Description)
		fmt.Printf("Created: %s\n", foundProject.CreatedAt.Format("2006-01-02 15:04:05"))

		if len(environments) == 0 {
			fmt.Println("\nNo environments found for this project")
			return
		}

		fmt.Println("\nEnvironments:")
		fmt.Println("=============")
		for _, e := range environments {
			fmt.Printf("- %s: %s\n", e.Name, e.Description)
		}
	},
}

func init() {
	// Import command flags
	importCmd.Flags().StringVar(&projectName, "project", "", "Project name (required)")
	importCmd.Flags().StringVar(&environmentName, "env", "development", "Environment name (default: development)")
	importCmd.MarkFlagRequired("project")

	// Export command flags
	exportCmd.Flags().StringVar(&projectName, "project", "", "Project name (required)")
	exportCmd.Flags().StringVar(&environmentName, "env", "development", "Environment name (default: development)")
	exportCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwriting the file if it exists")
	exportCmd.MarkFlagRequired("project")

	// Set env command flags
	setEnvCmd.Flags().StringVar(&projectName, "project", "", "Project name (required)")
	setEnvCmd.Flags().StringVar(&environmentName, "env", "development", "Environment name (default: development)")
	setEnvCmd.Flags().StringVar(&keyName, "key", "", "Environment variable key (required)")
	setEnvCmd.Flags().StringVar(&keyValue, "value", "", "Environment variable value")
	setEnvCmd.MarkFlagRequired("project")
	setEnvCmd.MarkFlagRequired("key")

	// Get env command flags
	getEnvCmd.Flags().StringVar(&projectName, "project", "", "Project name (required)")
	getEnvCmd.Flags().StringVar(&environmentName, "env", "development", "Environment name (default: development)")
	getEnvCmd.Flags().StringVar(&keyName, "key", "", "Environment variable key (required)")
	getEnvCmd.MarkFlagRequired("project")
	getEnvCmd.MarkFlagRequired("key")

	// Delete env command flags
	deleteEnvCmd.Flags().StringVar(&projectName, "project", "", "Project name (required)")
	deleteEnvCmd.Flags().StringVar(&environmentName, "env", "development", "Environment name (default: development)")
	deleteEnvCmd.Flags().StringVar(&keyName, "key", "", "Environment variable key (required)")
	deleteEnvCmd.MarkFlagRequired("project")
	deleteEnvCmd.MarkFlagRequired("key")

	// List env command flags
	listEnvCmd.Flags().StringVar(&projectName, "project", "", "Project name (required)")
	listEnvCmd.Flags().StringVar(&environmentName, "env", "development", "Environment name (default: development)")
	listEnvCmd.Flags().StringVar(&keyName, "filter", "", "Filter by key pattern")
	listEnvCmd.MarkFlagRequired("project")

	// Delete project command flags
	softDeleteProjectCmd.Flags().StringVar(&projectName, "project", "", "Project name (required)")
	softDeleteProjectCmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")
	softDeleteProjectCmd.MarkFlagRequired("project")

	// Create environment command flags
	createEnvironmentCmd.Flags().StringVar(&environmentName, "name", "", "Environment name (required)")
	createEnvironmentCmd.Flags().StringVar(&description, "description", "", "Environment description")
	createEnvironmentCmd.MarkFlagRequired("name")

	// Project details command flags
	projectDetailsCmd.Flags().StringVar(&projectName, "project", "", "Project name (required)")
	projectDetailsCmd.MarkFlagRequired("project")

	// Add environment subcommands
	environmentCmd.AddCommand(listEnvironmentsCmd)
	environmentCmd.AddCommand(createEnvironmentCmd)
}
