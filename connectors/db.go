package connectors

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/hellodhlyn/luppiter/internal/env"
)

func NewDatabaseConnection() (*gorm.DB, error) {
	url := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.DatabaseHost, env.DatabasePort, env.DatabaseUsername, env.DatabasePassword, env.DatabaseName,
	)
	return gorm.Open(postgres.New(postgres.Config{DSN: url}), &gorm.Config{})
}
