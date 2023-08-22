package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
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
