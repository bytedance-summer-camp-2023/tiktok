package db

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

type Comment struct {
	ID          uint      `gorm:"primarykey"`
	CreatedTime time.Time `gorm:"index;not null" json:"create_date"`
	UpdatedTime time.Time
	DeletedTime time.Time `gorm:"index"`
	Video       Video     `gorm:"foreignkey:VideoID" json:"video,omitempty"`
	VideoID     uint      `gorm:"index:idx_videoid;not null" json:"video_id"`
	User        User      `gorm:"foreignkey:UserID" json:"user,omitempty"`
	UserID      uint      `gorm:"index:idx_userid;not null" json:"user_id"`
	Content     string    `gorm:"type:varchar(255);not null" json:"content"`
}

func (Comment) TableName() string {
	return "Comments"
}

func CreateComment(ctx context.Context, comment *Comment) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// add new comment
		err := tx.Create(comment).Error
		if err != nil {
			return err
		}

		// increase the comment number of the relative video
		res := tx.Model(&Video{}).Where("id = ?", comment.VideoID).
			Update("comment_count", gorm.Expr("comment_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		return nil
	})

	return err
}

func DelCommentByID(ctx context.Context, commentID int64, vid int64) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		comment := new(Comment)

		// find comment in database
		err := tx.First(&comment, commentID).Error
		if err != nil {
			return err
		} else if err == gorm.ErrRecordNotFound {
			return nil
		}

		// delete comment
		err = tx.Where("id = ?", commentID).Delete(&Comment{}).Error
		if err != nil {
			return err
		}

		// decrease comment number of the relative video
		res := tx.Model(&Video{}).Where("id = ?", comment.VideoID).
			Update("comment_count", gorm.Expr("comment_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		return nil
	})

	return err
}

func GetVideoCommentListByVideoID(ctx context.Context, videoID int64) ([]*Comment, error) {
	var comments []*Comment

	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Model(&Comment{}).
		Where(&Comment{VideoID: uint(videoID)}).Order("created_time DESC").Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func GetCommentByCommentID(ctx context.Context, commentID int64) (*Comment, error) {
	comment := new(Comment)

	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Where("id = ?", commentID).First(&comment).Error
	if err == nil {
		return comment, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}
