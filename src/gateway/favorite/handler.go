package favorite

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"tiktok/src/constant/config"
	"tiktok/src/constant/strings"
	"tiktok/src/extra/tracing"
	"tiktok/src/gateway/models"
	"tiktok/src/rpc/favorite"
	grpc2 "tiktok/src/utils/grpc"
	"tiktok/src/utils/logging"
)

var Client favorite.FavoriteServiceClient

func init() {
	conn := grpc2.Connect(config.FavoriteRpcServerName)
	Client = favorite.NewFavoriteServiceClient(conn)
}

func ActionFavoriteHandler(c *gin.Context) {
	var req models.ActionFavoriteReq
	_, span := tracing.Tracer.Start(c.Request.Context(), "ActionFavoriteHandler")
	defer span.End()
	logger := logging.LogService("GateWay.ActionFavorite").WithContext(c.Request.Context())

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, models.ActionCommentRes{
			StatusCode: strings.GateWayParamsErrorCode,
			StatusMsg:  strings.GateWayParamsError,
		})
		return
	}

	actionType := uint32(req.ActionType)
	res, err := Client.FavoriteAction(c.Request.Context(), &favorite.FavoriteRequest{
		ActorId:    uint32(req.ActorId),
		VideoId:    uint32(req.VideoId),
		ActionType: actionType,
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"ActorId":    req.ActorId,
			"VideoId":    req.VideoId,
			"ActionType": req.ActionType,
		}).Warnf("Error when trying to connect with ActionFavoriteService")
		c.JSON(http.StatusOK, res)
		return
	}

	logger.WithFields(logrus.Fields{
		"ActorId":    req.ActorId,
		"VideoId":    req.VideoId,
		"ActionType": req.ActionType,
	}).Infof("Action favorite success")

	c.JSON(http.StatusOK, res)
}

func ListFavoriteHandler(c *gin.Context) {
	var req models.ListFavoriteReq
	_, span := tracing.Tracer.Start(c.Request.Context(), "ListFavoriteHandler")
	defer span.End()
	logger := logging.LogService("GateWay.ListFavorite").WithContext(c.Request.Context())

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, models.ListCommentRes{
			StatusCode: strings.GateWayParamsErrorCode,
			StatusMsg:  strings.GateWayParamsError,
		})
		return
	}

	res, err := Client.FavoriteList(c.Request.Context(), &favorite.FavoriteListRequest{
		ActorId: uint32(req.ActorId),
		UserId:  uint32(req.UserId),
	})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ActorId": req.ActorId,
			"UserId":  req.UserId,
		}).Warnf("Error when trying to connect with ListFavoriteHandler")
		c.JSON(http.StatusOK, res)
		return
	}

	logger.WithFields(logrus.Fields{
		"ActorId": req.ActorId,
		"UserId":  req.UserId,
	}).Infof("List favorite videos success")

	c.JSON(http.StatusOK, res)
}
