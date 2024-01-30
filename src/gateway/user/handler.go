package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"tiktok/src/constant/config"
	"tiktok/src/constant/strings"
	"tiktok/src/extra/tracing"
	"tiktok/src/gateway/models"
	"tiktok/src/rpc/user"
	grpc2 "tiktok/src/utils/grpc"
	"tiktok/src/utils/logging"
)

var userClient user.UserServiceClient

func init() {
	userConn := grpc2.Connect(config.UserRpcServerName)
	userClient = user.NewUserServiceClient(userConn)
}

func UserHandler(c *gin.Context) {
	var req models.UserReq
	_, span := tracing.Tracer.Start(c.Request.Context(), "UserInfoHandler")
	defer span.End()
	logger := logging.LogService("GateWay.UserInfo").WithContext(c.Request.Context())

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, models.UserRes{
			StatusCode: strings.GateWayParamsErrorCode,
			StatusMsg:  strings.GateWayParamsError,
		})
		logging.SetSpanError(span, err)
		return
	}

	resp, err := userClient.GetUserInfo(c.Request.Context(), &user.UserRequest{
		UserId:  req.UserId,
		ActorId: req.ActorId,
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when gateway get info from UserInfo Service")
		logging.SetSpanError(span, err)
		c.JSON(http.StatusOK, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
