package database

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func getenvOrDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func init() {
	newDB, err := gorm.Open("postgres", "host="+getenvOrDefault("DB_HOST", "127.0.0.1")+" port="+getenvOrDefault("DB_PORT", "5432")+" user="+os.Getenv("DB_USERNAME")+" password="+os.Getenv("DB_PASSWORD")+" dbname="+os.Getenv("DB_NAME")+" sslmode=disable")
	if err != nil {
		panic(err)
	}

	DB = newDB
}
