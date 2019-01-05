package services

import (
	"context"
	"testing"

	"luppiter/components/database"
)

var (
	testCtx           context.Context
	testUserUUID      = "example_user_uuid"
	testAPIKey        = "example_api_key"
	testService       = "Test"
	testPermission    = "Test::*"
	testAuthorization = "example_authorization"
)

func TestMain(m *testing.M) {
	database.DB.AutoMigrate(
		APIKey{},
		KeyValueItem{},
		StorageBucket{},
		StorageItem{},
	)

	testCtx = context.Background()
	testCtx = context.WithValue(testCtx, "APIKey", testAPIKey)

	database.DB.Save(&APIKey{UserUUID: testUserUUID, Comment: "-", Key: testAPIKey, Permissions: []string{testPermission}})

	m.Run()
}
