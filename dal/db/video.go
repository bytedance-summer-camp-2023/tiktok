package db

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

type Video struct {
	ID            uint      `gorm:"primarykey"`
	CreatedAt     time.Time `gorm:"not null;index:idx_create" json:"created_at,omitempty"`
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Author        User           `gorm:"foreignkey:AuthorID" json:"author,omitempty"`
	AuthorID      uint           `gorm:"index:idx_authorid;not null" json:"author_id,omitempty"`
	PlayUrl       string         `gorm:"type:varchar(255);not null" json:"play_url,omitempty"`
	CoverUrl      string         `gorm:"type:varchar(255)" json:"cover_url,omitempty"`
	FavoriteCount uint           `gorm:"default:0;not null" json:"favorite_count,omitempty"`
	CommentCount  uint           `gorm:"default:0;not null" json:"comment_count,omitempty"`
	Title         string         `gorm:"type:varchar(50);not null" json:"title,omitempty"`
}

func (Video) TableName() string {
	return "videos"
}

func GetVideosList(ctx context.Context, limit int, latestTime *int64) ([]*Video, error) {
	// 初始化视频切片
	videos := make([]*Video, 0)

	// 判断 latestTime 是否为空或零值，如果是则设置为当前时间
	if latestTime == nil || *latestTime == 0 {
		curTime := time.Now().UnixMilli()
		latestTime = &curTime
	}

	// 获取数据库连接
	conn := GetDB().Clauses(dbresolver.Read).WithContext(ctx)

	// 构造查询条件
	query := conn.Limit(limit).Order("created_at desc").Where("created_at < ?", time.UnixMilli(*latestTime))

	// 执行查询并处理错误
	if err := query.Find(&videos).Error; err != nil {
		return nil, err
	}

	// 成功，返回视频切片
	return videos, nil
}
