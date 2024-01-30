package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
	"strconv"
	"tiktok/src/constant/config"
	"tiktok/src/constant/strings"
	"tiktok/src/extra/tracing"
	"tiktok/src/rpc/auth"
	grpc2 "tiktok/src/utils/grpc"
	"tiktok/src/utils/logging"
)

var client auth.AuthServiceClient

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/douyin/user/login" || c.Request.URL.Path == "/douyin/user/register" {
			c.Next()
			return
		}

		var token string
		if c.Request.URL.Path == "/douyin/publish/action/" {
			token = c.PostForm("token")
		} else {
			token = c.Query("token")
		}
		ctx, span := tracing.Tracer.Start(c.Request.Context(), "AuthMiddleWare")
		defer span.End()
		span.SetAttributes(attribute.String("token", token))
		logger := logging.LogService("GateWay.AuthMiddleWare").WithContext(ctx)
		authenticate, err := client.Authenticate(c.Request.Context(), &auth.AuthenticateRequest{Token: token})

		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Gatewat Auth meet trouble")

			c.JSON(http.StatusOK, gin.H{
				"status_code": strings.GateWayErrorCode,
				"status_msg":  strings.GateWayError,
			})
			c.Abort()
			return
		}

		if authenticate.StatusCode != 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": strings.AuthUserNeededCode,
				"status_msg":  strings.AuthUserNeeded,
			})
			c.Abort()
			return
		}
		c.Request.URL.RawQuery += "&actor_id=" + strconv.FormatUint(uint64(authenticate.UserId), 10)
		c.Next()

	}
}

func init() {
	authConn := grpc2.Connect(config.AuthRpcServerName)
	client = auth.NewAuthServiceClient(authConn)
}
