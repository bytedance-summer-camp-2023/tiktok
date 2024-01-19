package models

import (
	"gorm.io/gorm"
	database "tiktok/src/storage/db"
)

type Comment struct {
	ID      uint32 `gorm:"not null;index:comment_video;primarykey"`                        // 评论 ID
	VideoId uint32 `json:"video_id" column:"video_id" gorm:"not null;index:comment_video"` // 视频 ID
	UserId  uint32 `json:"user_id" column:"user_id" gorm:"not null"`                       // 用户 ID
	Content string `json:"content" column:"content"`                                       // 评论内容
	gorm.Model
}

func init() {
	if err := database.Client.AutoMigrate(&Comment{}); err != nil {
		panic(err)
	}
}
