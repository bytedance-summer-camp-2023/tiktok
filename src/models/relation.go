package models

import (
	"gorm.io/gorm"
	database "tiktok/src/storage/db"
)

type Relation struct {
	ID      uint32 `gorm:"not null;index:relation;primarykey;autoIncrement"`            // relation ID
	ActorId uint32 `json:"actor_id" column:"actor_id" gorm:"not null;index:actor_list"` // 粉丝 ID
	UserId  uint32 `json:"user_id" column:"user_id" gorm:"not null;index:user_list"`    // 被关注用户 ID
	gorm.Model
}

func init() {
	if err := database.Client.AutoMigrate(&Relation{}); err != nil {
		panic(err)
	}
}