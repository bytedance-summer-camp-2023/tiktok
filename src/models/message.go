package models

import (
	"gorm.io/gorm"
	database "tiktok/src/storage/db"
)

type Message struct {
	ID           uint32 `gorm:"not null;primarykey;autoIncrement"`
	To_user_id   uint32 `gorm:"not null"`
	From_user_id uint32 `gorm:"not null"`
	Content      string `gorm:"not null"`
	// Create_time  time.Time `gorm:"not null"`
	gorm.Model
}

func init() {
	if err := database.Client.AutoMigrate(&Message{}); err != nil {
		panic(err)
	}
}
