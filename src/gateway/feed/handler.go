package feed

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"tiktok/src/constant/config"
	"tiktok/src/constant/strings"
	"tiktok/src/extra/tracing"
	"tiktok/src/gateway/models"
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
		c.JSON(http.StatusOK, models.ListVideosRes{
			StatusCode: strings.GateWayParamsErrorCode,
			StatusMsg:  strings.GateWayParamsError,
			NextTime:   nil,
			VideoList:  nil,
		})
		return
	}

	res, err := Client.ListVideos(c.Request.Context(), &feed.ListFeedRequest{})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"LatestTime": req.LatestTime,
		}).Warnf("Error when trying to connect with FeedService")
		c.JSON(http.StatusOK, models.ListVideosRes{
			StatusCode: strings.FeedServiceInnerErrorCode,
			StatusMsg:  strings.FeedServiceInnerError,
			NextTime:   nil,
			VideoList:  nil,
		})
		return
	}
	latestTime := req.LatestTime
	if _, err := strconv.Atoi(latestTime); latestTime != "" && err != nil {
		logger.WithFields(logrus.Fields{
			"LatestTime": req.LatestTime,
		}).Warnf("Error when trying to convert LatestTime to int")
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
	}).Infof("Feed List videos")
	c.JSON(http.StatusOK, res)
}

func init() {
	conn, err := grpc2.Connect(config.FeedRpcServerName)
	if err != nil {
		panic(err)
	}
	Client = feed.NewFeedServiceClient(conn)
}
