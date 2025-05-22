package utils

import (
	"fmt"
	"os"
	"strings"
)

// ParseKeyValuePair parses a string in the format "key=value"
func ParseKeyValuePair(input string) (string, string, error) {
	parts := strings.SplitN(input, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format: %s (expected key=value)", input)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	return key, value, nil
}

// GetUserConfirmation asks for user confirmation (y/n)
func GetUserConfirmation(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

// EnsureFileExists checks if a file exists
func EnsureFileExists(filePath string) error {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}
	if info.IsDir() {
		return fmt.Errorf("expected file, got directory: %s", filePath)
	}
	return nil
}

// FormatEnvValue formats a value for env file output
func FormatEnvValue(value string) string {
	// If value contains spaces or special characters, quote it
	if strings.ContainsAny(value, " \t\n\r\"'`) ") {
		return fmt.Sprintf("\"%s\"", value)
	}
	return value
}
