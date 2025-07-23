package database

import (
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
	//&user.User{},
	//&chatroom.ChatRoom{},
	//&message.Message{},
	)
}

func MigrateTable(db *gorm.DB, model interface{}) error {
	return db.AutoMigrate(model)
}
