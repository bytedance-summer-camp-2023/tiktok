package service

import (
	"context"
	"time"

	"github.com/bytedance-summer-camp-2023/tiktok/dal/db"
	user "github.com/bytedance-summer-camp-2023/tiktok/kitex/kitex_gen/user"
	video "github.com/bytedance-summer-camp-2023/tiktok/kitex/kitex_gen/video"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/minio"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/zap"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

const limit = 30

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	logger := zap.InitLogger()
	nextTime := time.Now().UnixMilli()
	var userID int64 = -1

	// 验证token有效性
	if req.Token != "" {
		claims, err := Jwt.ParseToken(req.Token)
		if err != nil {
			logger.Errorln(err.Error())
			res := &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "token 解析错误",
			}
			return res, nil
		}
		userID = claims.Id
	}
	// 调用数据库查询 video_list
	videos, err := db.MGetVideos(ctx, limit, &req.LatestTime)
	if err != nil {
		logger.Errorln(err.Error())
		res := &video.FeedResponse{
			StatusCode: -1,
			StatusMsg:  "视频获取失败：服务器内部错误",
		}
		return res, nil
	}
	videoList := make([]*video.Video, 0)
	for _, r := range videos {
		author, err := db.GetUserByID(ctx, int64(r.AuthorID))
		if err != nil {
			logger.Errorf("error:%v", err.Error())
			return nil, err
		}
		relation, err := db.GetRelationByUserIDs(ctx, userID, int64(author.ID))
		if err != nil {
			logger.Errorln(err.Error())
			res := &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "视频获取失败：服务器内部错误",
			}
			return res, nil
		}
		favorite, err := db.GetFavoriteVideoRelationByUserVideoID(ctx, userID, int64(r.ID))
		if err != nil {
			logger.Errorln(err.Error())
			res := &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "视频获取失败：服务器内部错误",
			}
			return res, nil
		}
		playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, r.PlayUrl)
		if err != nil {
			logger.Errorf("Minio获取链接失败：%v", err.Error())
			res := &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误：视频获取失败",
			}
			return res, nil
		}
		coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, r.CoverUrl)
		if err != nil {
			logger.Errorf("Minio获取链接失败：%v", err.Error())
			res := &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误：封面获取失败",
			}
			return res, nil
		}
		avatarUrl, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, author.Avatar)
		if err != nil {
			logger.Errorf("Minio获取链接失败：%v", err.Error())
			res := &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误：头像获取失败",
			}
			return res, nil
		}
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, author.BackgroundImage)
		if err != nil {
			logger.Errorf("Minio获取链接失败：%v", err.Error())
			res := &video.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误：背景图获取失败",
			}
			return res, nil
		}

		videoList = append(videoList, &video.Video{
			Id: int64(r.ID),
			Author: &user.User{
				Id:              int64(author.ID),
				Name:            author.UserName,
				FollowCount:     int64(author.FollowingCount),
				FollowerCount:   int64(author.FollowerCount),
				IsFollow:        relation != nil,
				Avatar:          avatarUrl,
				BackgroundImage: backgroundUrl,
				Signature:       author.Signature,
				TotalFavorited:  int64(author.TotalFavorited),
				WorkCount:       int64(author.WorkCount),
				FavoriteCount:   int64(author.FavoriteCount),
			},
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: int64(r.FavoriteCount),
			CommentCount:  int64(r.CommentCount),
			IsFavorite:    favorite != nil,
			Title:         r.Title,
		})
	}
	if len(videos) != 0 {
		nextTime = videos[len(videos)-1].UpdatedAt.UnixMilli()
	}
	res := &video.FeedResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
		NextTime:   nextTime,
	}
	return res, nil
}

// PublishAction implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {

	return nil, nil
}

// PublishList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	return nil, nil
}
