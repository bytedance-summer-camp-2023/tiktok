package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"tiktok/src/constant/config"
	"tiktok/src/constant/strings"
	"tiktok/src/gateway/models"
	"tiktok/src/rpc/auth"
	"tiktok/src/utils/interceptor"
	"tiktok/src/utils/logging"
	"tiktok/src/utils/trace"
)

var Client auth.AuthServiceClient

func LoginHandle(c *gin.Context) {
	var req models.LoginReq
	span := trace.GetChildSpanFromGinContext(c, "GateWay-Login")
	defer span.Finish()
	log := logging.GetSpanLogger(span, "GateWay.Login")

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, models.LoginRes{
			StatusCode: strings.GateWayParamsErrorCode,
			StatusMsg:  strings.GateWayParamsError,
			UserId:     0,
			Token:      "",
		})
	}

	res, err := Client.Login(c.Request.Context(), &auth.LoginRequest{
		Username: req.UserName,
		Password: req.Password,
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"Username": req.UserName,
		}).Warnf("Error when trying to connect with AuthService")
		c.JSON(http.StatusOK, res)
		return
	}

	log.WithFields(logrus.Fields{
		"Username": req.UserName,
		"Token":    res.Token,
		"UserId":   res.UserId,
	}).Infof("User log in")

	c.JSON(http.StatusOK, res)
}

func RegisterHandle(c *gin.Context) {
	var req models.RegisterReq
	span := trace.GetChildSpanFromGinContext(c, "GateWay-Register")
	defer span.Finish()
	log := logging.GetSpanLogger(span, "GateWay.Register")

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, models.LoginRes{
			StatusCode: strings.GateWayParamsErrorCode,
			StatusMsg:  strings.GateWayParamsError,
			UserId:     0,
			Token:      "",
		})
	}

	res, err := Client.Register(c.Request.Context(), &auth.RegisterRequest{
		Username: req.UserName,
		Password: req.Password,
	})

	if err != nil {
		log.WithFields(logrus.Fields{
			"Username": req.UserName,
		}).Warnf("Error when trying to connect with AuthService")
		c.JSON(http.StatusOK, res)
		return
	}

	log.WithFields(logrus.Fields{
		"Username": req.UserName,
		"Token":    res.Token,
		"UserId":   res.UserId,
	}).Infof("User register in")

	c.JSON(http.StatusOK, res)
}

func init() {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s/%s?wait=15s", config.EnvCfg.ConsulAddr, config.AuthRpcServerName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(interceptor.OpenTracingClientInterceptor()),
	)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Build AuthService Cient meet trouble")
	}
	Client = auth.NewAuthServiceClient(conn)
}
