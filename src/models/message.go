package models

import (
	"gorm.io/gorm"
	database "tiktok/src/storage/db"
)

type Message struct {
	ID             uint32 `gorm:"not null;primarykey;autoIncrement"`
	ToUserId       uint32 `gorm:"not null" `
	FromUserId     uint32 `gorm:"not null"`
	ConversationId string `gorm:"not null" index:"conversationid"`
	Content        string `gorm:"not null"`
	gorm.Model
}

func init() {
	if err := database.Client.AutoMigrate(&Message{}); err != nil {
		panic(err)
	}
}
