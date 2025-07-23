package database

import (
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
	//&entity.User{},
	// Add other models here as needed
	)
}

func MigrateTable(db *gorm.DB, model interface{}) error {
	return db.AutoMigrate(model)
}
