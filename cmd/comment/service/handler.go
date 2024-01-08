package service

import (
	"context"
	"gorm.io/gorm"
	"tiktok/dal/db"
	comment "tiktok/kitex/kitex_gen/comment"
	"tiktok/kitex/kitex_gen/user"
	"tiktok/pkg/minio"
	"tiktok/pkg/zap"
	"time"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CommentAction implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentAction(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	logger := zap.InitLogger()

	// user authentication
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		logger.Errorf("Token parsing error: %v", err.Error())
		res := &comment.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "Token parsing error",
		}
		return res, nil
	}

	// get video
	userID := claims.Id
	video, _ := db.GetVideoById(ctx, req.VideoId)
	if video == nil {
		logger.Errorf("Video ID doesn't exitst: %v", req.VideoId)
		res := &comment.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "Video does not exist: server inner error",
		}
		return res, nil
	}

	// comment action
	actionType := req.ActionType
	if actionType == 1 {
		// add comment with ActionType = 1
		cmt := &db.Comment{
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
			DeletedTime: time.Now(),
			VideoID:     uint(req.VideoId),
			UserID:      uint(userID),
			Content:     req.CommentText,
		}

		// create error
		err := db.CreateComment(ctx, cmt)
		if err != nil {
			logger.Errorf("Create comment failed：%v", err.Error())
			res := &comment.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "Create comment failed: server inner error",
			}
			return res, nil
		}
	} else if actionType == 2 {
		// delete comment with ActionType = 2
		cmt, err := db.GetCommentByCommentID(ctx, req.CommentId)
		if err != nil {
			logger.Errorf("Error occured while deleting comment：%v", err.Error())
			res := &comment.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "Delete comment failed: server inner error",
			}
			return res, nil
		}

		if cmt == nil {
			// comment does not exist error
			logger.Errorf("Comment doesn't exist：%v", req.CommentId)
			res := &comment.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "Delete comment failed: server inner error",
			}
			return res, nil
		} else {
			v, err := db.GetVideoById(ctx, int64(cmt.VideoID))
			if err != nil {
				logger.Errorf("Video doesn't exist：%v", err.Error())
				res := &comment.CommentActionResponse{
					StatusCode: -1,
					StatusMsg:  "Delete comment failed: server inner error",
				}
				return res, nil
			}

			// authentication error
			if userID != int64(cmt.UserID) || userID != int64(v.AuthorID) {
				logger.Errorf("Authentication failed with user ID：%v", cmt.UserID)
				res := &comment.CommentActionResponse{
					StatusCode: -1,
					StatusMsg:  "Delete comment failed: server inner error",
				}
				return res, nil
			}

			err = db.DelCommentByID(ctx, req.CommentId, req.VideoId)
			if err != nil {
				logger.Errorf("Delete comment failed：%v", err.Error())
				res := &comment.CommentActionResponse{
					StatusCode: -1,
					StatusMsg:  "Delete comment failed: server inner error",
				}
				return res, nil
			}
		}
	} else {
		// invalid ActionType
		res := &comment.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "Action_type invalid",
		}
		return res, nil
	}

	// success response
	res := &comment.CommentActionResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
	}

	return res, nil
}

// CommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	logger := zap.InitLogger()

	// user authentication
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		logger.Errorf("Token parsing error: %v", err.Error())
		res := &comment.CommentListResponse{
			StatusCode: -1,
			StatusMsg:  "Token parsing error",
		}
		return res, nil
	}
	userID := claims.Id

	// get comments from database
	results, err := db.GetVideoCommentListByVideoID(ctx, req.VideoId)
	if err != nil {
		logger.Errorf("Get comment list failed: %v", err.Error())
		res := &comment.CommentListResponse{
			StatusCode: -1,
			StatusMsg:  "Get comment list failed: server inner error",
		}
		return res, nil
	}

	// fill the comment list
	comments := make([]*comment.Comment, 0)
	for _, r := range results {
		u, err := db.GetUserByID(ctx, int64(r.UserID))
		if err != nil {
			logger.Errorf("Get user ID failed: %v", err.Error())
			res := &comment.CommentListResponse{
				StatusCode: -1,
				StatusMsg:  "Get user ID failed：server inner error",
			}
			return res, nil
		}

		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, u.Avatar)
		if err != nil {
			logger.Errorf("Minio error while getting avatar URL：%v", err.Error())
			res := &comment.CommentListResponse{
				StatusCode: -1,
				StatusMsg:  "Get avatar failed: server inner error",
			}
			return res, nil
		}

		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, u.Avatar)
		if err != nil {
			logger.Errorf("Minio error while getting background URL：%v", err.Error())
			res := &comment.CommentListResponse{
				StatusCode: -1,
				StatusMsg:  "Get background failed: server inner error",
			}
			return res, nil
		}

		usr := &user.User{
			Id:              userID,
			Name:            u.UserName,
			FollowCount:     int64(u.FollowingCount),
			FollowerCount:   int64(u.FollowerCount),
			IsFollow:        err != gorm.ErrRecordNotFound,
			Avatar:          avatar,
			BackgroundImage: backgroundUrl,
			Signature:       u.Signature,
			TotalFavorited:  int64(u.TotalFavorited),
			WorkCount:       int64(u.WorkCount),
			FavoriteCount:   int64(u.FavoriteCount),
		}
		comments = append(comments, &comment.Comment{
			Id:         int64(r.ID),
			User:       usr,
			Content:    r.Content,
			CreateDate: r.CreatedTime.Format("2006-01-02"),
		})
	}

	// success response
	res := &comment.CommentListResponse{
		StatusCode:  0,
		StatusMsg:   "Success",
		CommentList: comments,
	}

	return res, nil
}
