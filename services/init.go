package services

import "luppiter/components/database"

func init() {
	database.DB.AutoMigrate(
		APIKey{},
		KeyValueItem{},
		StorageBucket{},
		StorageItem{},
	)
}
