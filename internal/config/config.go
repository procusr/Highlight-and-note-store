package config

import "os"

// Config holds all configuration for the application
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	Environment string
	OpenAIKey   string
}

// Load loads configuration from environment variables with sensible defaults
func Load() *Config {
	cfg := &Config{
		Port:        getEnvWithDefault("PORT", "8080"),
		DatabaseURL: getEnvWithDefault("DATABASE_URL", "postgres://user:password@localhost:5433/mhgo?sslmode=disable"),
		JWTSecret:   getEnvWithDefault("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		Environment: getEnvWithDefault("ENVIRONMENT", "development"),
		OpenAIKey:   getEnvWithDefault("OPENAI_API_KEY", ""),
	}

	return cfg
}

// getEnvWithDefault returns the value of an environment variable or a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
