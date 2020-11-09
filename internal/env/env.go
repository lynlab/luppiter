package env

import "os"

var (
	DatabaseHost     = getEnvOrDefault("DB_HOST", "127.0.0.1")
	DatabasePort     = getEnvOrDefault("DB_PORT", "5432")
	DatabaseName     = getEnvOrDefault("DB_NAME", "luppiter")
	DatabaseUsername = getEnvOrDefault("DB_USERNAME", "postgres")
	DatabasePassword = getEnvOrDefault("DB_PASSWORD", "rootpass")
)

func getEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
