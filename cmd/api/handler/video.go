package handler

import (
	"bytes"
	"context"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/zap"
	"github.com/cloudwego/hertz/pkg/app"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/bytedance-summer-camp-2023/tiktok/cmd/api/rpc"
	"github.com/bytedance-summer-camp-2023/tiktok/internal/response"
	kitex "github.com/bytedance-summer-camp-2023/tiktok/kitex/kitex_gen/video"
)

func Feed(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	latestTime := c.Query("latest_time")
	var timestamp int64 = 0
	if latestTime != "" {
		timestamp, _ = strconv.ParseInt(latestTime, 10, 64)
	} else {
		timestamp = time.Now().UnixMilli()
	}

	req := &kitex.FeedRequest{
		LatestTime: timestamp,
		Token:      token,
	}
	res, _ := rpc.Feed(ctx, req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.Feed{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.Feed{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		VideoList: res.VideoList,
	})
}

func PublishAction(ctx context.Context, c *app.RequestContext) {
	logger := zap.InitLogger()
	token := c.PostForm("token")
	if token == "" {
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "用户鉴权失败，token为空",
			},
		})
		return
	}
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "标题不能为空",
			},
		})
		return
	}
	// 视频数据
	file, err := c.FormFile("data")
	if err != nil {
		logger.Errorln(err.Error())
		c.JSON(http.StatusBadRequest, response.RelationAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "上传视频加载失败",
			},
		})
		return
	}
	src, err := file.Open()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		logger.Errorln(err.Error())
		c.JSON(http.StatusBadRequest, response.RelationAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "视频上传失败",
			},
		})
		return
	}

	req := &kitex.PublishActionRequest{
		Token: token,
		Title: title,
		Data:  buf.Bytes(),
	}
	res, _ := rpc.PublishAction(ctx, req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.PublishAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.PublishAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
	})
}
