package models

import (
	"gorm.io/gorm"
	database "tiktok/src/storage/db"
)

// Video 视频表
type Video struct {
	ID        uint32 `gorm:"not null;index:video;primarykey;autoIncrement"`
	Title     string `json:"title" gorm:"not null;"`
	FileName  string `json:"play_name" gorm:"not null;"`
	CoverName string `json:"cover_name" gorm:"not null;"`
	gorm.Model
}

func init() {
	if err := database.Client.AutoMigrate(&Video{}); err != nil {
		panic(err)
	}
}