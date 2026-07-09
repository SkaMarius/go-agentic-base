package config

import "os"

const defaultPort = "8080"

// Config holds environment-derived application settings.
type Config struct {
	Port        string
	DatabaseURL string
}

// Load reads configuration from environment variables.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	return Config{
		Port:        port,
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
