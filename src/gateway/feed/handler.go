package feed

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"tiktok/src/constant/config"
	"tiktok/src/constant/strings"
	"tiktok/src/extra/tracing"
	"tiktok/src/gateway/models"
	"tiktok/src/gateway/utils"
	"tiktok/src/rpc/feed"
	grpc2 "tiktok/src/utils/grpc"
	"tiktok/src/utils/logging"
)

var Client feed.FeedServiceClient

func ListVideosHandle(c *gin.Context) {
	var req models.ListVideosReq
	_, span := tracing.Tracer.Start(c.Request.Context(), "Feed-ListVideoHandle")
	defer span.End()
	logger := logging.LogService("GateWay.Videos").WithContext(c.Request.Context())

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.WithFields(logrus.Fields{
			"latestTime": req.LatestTime,
			"err":        err,
		}).Warnf("Error when trying to bind query")
		c.JSON(http.StatusOK, models.ListVideosRes{
			StatusCode: strings.GateWayParamsErrorCode,
			StatusMsg:  strings.GateWayParamsError,
			NextTime:   nil,
			VideoList:  nil,
		})
	}

	latestTime := req.LatestTime
	res, err := Client.ListVideos(c.Request.Context(), &feed.ListFeedRequest{
		LatestTime: &latestTime,
	})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"LatestTime": latestTime,
		}).Warnf("Error when trying to connect with FeedService")
		c.JSON(http.StatusOK, models.ListVideosRes{
			StatusCode: strings.FeedServiceInnerErrorCode,
			StatusMsg:  strings.FeedServiceInnerError,
			NextTime:   nil,
			VideoList:  nil,
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"LatestTime": latestTime,
		"res":        res,
	}).Infof("Feed List videos")
	c.Render(http.StatusOK, utils.CustomJSON{Data: res, Context: c})
}

func init() {
	conn := grpc2.Connect(config.FeedRpcServerName)
	Client = feed.NewFeedServiceClient(conn)
}
