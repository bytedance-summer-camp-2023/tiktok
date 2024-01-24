package message

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"tiktok/src/constant/config"
	"tiktok/src/constant/strings"
	"tiktok/src/extra/tracing"
	"tiktok/src/gateway/models"
	"tiktok/src/rpc/chat"
	grpc2 "tiktok/src/utils/grpc"
	"tiktok/src/utils/logging"
)

var Client chat.ChatServiceClient

func init() {
	conn := grpc2.Connect(config.MessageRpcServerName)
	Client = chat.NewChatServiceClient(conn)
}

func ActionMessageHandler(c *gin.Context) {
	var req models.SMessageReq
	_, span := tracing.Tracer.Start(c.Request.Context(), "ActionMessageHandler")
	defer span.End()
	logger := logging.LogService("GateWay.ActionMessage").WithContext(c.Request.Context())

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.WithFields(logrus.Fields{
			//"CreateTime": req.Create_time,
			"err": err,
		}).Warnf("Error when trying to bind query")
		c.JSON(http.StatusOK, models.ActionCommentRes{
			StatusCode: strings.GateWayParamsErrorCode,
			StatusMsg:  strings.GateWayParamsError,
		})
		return
	}

	var res *chat.ActionResponse
	var err error

	res, err = Client.ChatAction(c.Request.Context(), &chat.ActionRequest{
		ActorId:    uint32(req.ActorId),
		UserId:     uint32(req.User_id),
		ActionType: uint32(req.Action_type),
		Content:    req.Content,
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"actor_id": req.ActorId,
			"content":  req.Content,
		}).Warnf("Error when trying to connect with ActionMessageHandler")

		//这个位置返回状态是不是有问题？
		c.JSON(http.StatusOK, res)
		return
	}
	logger.WithFields(logrus.Fields{
		"actor_id": req.ActorId,
		"content":  req.Content,
	}).Infof("Action send message success")

	c.JSON(http.StatusOK, res)
}

func ListMessageHandler(c *gin.Context) {
	var req models.ListMessageReq
	_, span := tracing.Tracer.Start(c.Request.Context(), "ListMessageHandler")
	defer span.End()
	logger := logging.LogService("GateWay.ListMessage").WithContext(c.Request.Context())

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, models.ListCommentRes{
			StatusCode: strings.GateWayParamsErrorCode,
			StatusMsg:  strings.GateWayParamsError,
		})
		return
	}

	res, err := Client.Chat(c.Request.Context(), &chat.ChatRequest{
		ActorId:    req.ActorId,
		UserId:     req.UserId,
		PreMsgTime: req.PreMsgTime,
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"actor_id": req.ActorId,
			"user_id":  req.UserId,
		}).Warnf("Error when trying to connect with ListMessageHandler")
		c.JSON(http.StatusOK, res)
		return
	}

	logger.WithFields(logrus.Fields{
		"actor_id": req.ActorId,
		"user_id":  req.UserId,
	}).Infof("List comment success")

	c.JSON(http.StatusOK, res)
}
