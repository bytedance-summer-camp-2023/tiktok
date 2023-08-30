package db

import (
	"context"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/errno"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func CreateVideo(ctx context.Context, video *Video) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 在 video 表中创建视频记录
		err := tx.Create(video).Error
		if err != nil {
			return err
		}
		// 2. 同步 user 表中的作品数量
		res := tx.Model(&User{}).Where("id = ?", video.AuthorID).Update("work_count", gorm.Expr("work_count + ?", 1))
		if res.Error != nil {
			return err
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}
		return nil
	})

	return err
}

func GetVideosByUserID(ctx context.Context, authorId int64) ([]*Video, error) {
	var publishList []*Video
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Model(&Video{}).
		Where(&Video{AuthorID: uint(authorId)}).Find(&publishList).Error
	if err != nil {
		return nil, err
	}

	return publishList, nil
}

func DelVideoByID(ctx context.Context, videoID int64, authorID int64) error {
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 根据主键 video_id 删除 video
		err := tx.Unscoped().Delete(&Video{}, videoID).Error
		if err != nil {
			return err
		}
		// 2. 同步 user 表中的作品数量
		res := tx.Model(&User{}).Where("id = ?", authorID).Update("work_count", gorm.Expr("work_count - ?", 1))
		if res.Error != nil {
			return err
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}
		return nil
	})
	return err
}
