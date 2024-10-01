package config

import (
	"log"
	"os"
)

// GetEnv retrieves the value of the environment variable identified by the key.
// If the environment variable is not set, it logs a message and returns the provided defaultValue.
//
// Parameters:
//   - key: The name of the environment variable to retrieve.
//   - defaultValue: The value to return if the environment variable is not set.
//
// Returns:
//   - A string containing the value of the environment variable, or defaultValue if it is not set.
func GetEnv(key string, defaultValue string) string {
	// Retrieve the value of the environment variable
	value := os.Getenv(key)

	// Check if the value is empty, indicating the environment variable is not set
	if value == "" {
		// Log a message indicating that the default value is being used
		log.Printf("Using default value for %s: %s", key, defaultValue)
		// Return the default value
		return defaultValue
	}

	// Return the value of the environment variable
	return value
}
