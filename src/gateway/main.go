package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"tiktok/src/constant/config"
	"tiktok/src/extra/profiling"
	"tiktok/src/extra/tracing"
	"tiktok/src/gateway/about"
	"tiktok/src/gateway/auth"
	comment2 "tiktok/src/gateway/comment"
	feed2 "tiktok/src/gateway/feed"
	message2 "tiktok/src/gateway/message"
	"tiktok/src/gateway/middleware"
	relation2 "tiktok/src/gateway/relation"
	"tiktok/src/utils/logging"
)

func main() {
	// Set Trace Provider
	tp, err := tracing.SetTraceProvider(config.WebServiceName)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panicf("Error to set the trace")
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error to set the trace")
		}
	}()
	g := gin.Default()
	// Configure Gzip
	g.Use(gzip.Gzip(gzip.DefaultCompression))
	// Configure Tracing
	g.Use(otelgin.Middleware(config.WebServiceName))
	g.Use(middleware.TokenAuthMiddleware())

	// Configure Pyroscope
	profiling.InitPyroscope("TikTok.GateWay")
	// Configure Prometheus
	p := ginprometheus.NewPrometheus("gin")
	p.Use(g)

	// Register Service
	// Test Service
	g.GET("/about", about.Handle)

	// Production Service
	rootPath := g.Group("/douyin")
	user := rootPath.Group("/user")
	{
		user.POST("/login", auth.LoginHandle)
		user.POST("/register", auth.RegisterHandle)
	}
	feed := rootPath.Group("/feed")
	{
		feed.GET("/", feed2.ListVideosHandle)
	}
	comment := rootPath.Group("/comment")
	{
		comment.POST("/action", comment2.ActionCommentHandler)
		comment.GET("/list", comment2.ListCommentHandler)
		comment.GET("/count", comment2.CountCommentHandler)
	}
	//todo
	message := rootPath.Group("/message")
	{
		message.GET("/chat", message2.ListMessageHandler)
		message.POST("/action", message2.ActionMessageHandler)
		fmt.Println(message)
	}

	relation := rootPath.Group("/relation")
	{
		//todo: frontend
		//relation.POST("/action", relation2.ActionHandler)
		relation.POST("/follow", relation2.FollowHandler)
		relation.POST("/unfollow", relation2.UnfollowHandler)
		relation.GET("/follow/list", relation2.GetFollowListHandler)
		relation.GET("/follower/list", relation2.GetFollowerListHandler)
		relation.GET("/friend/list", relation2.GetFriendListHandler)
		relation.GET("/follow/count", relation2.CountFollowHandler)
		relation.GET("/follower/count", relation2.CountFollowerHandler)
		relation.GET("/isFollow", relation2.IsFollowHandler)
	}
	// Run Server
	if err := g.Run(config.WebServiceAddr); err != nil {
		panic("Can not run TikTok Gateway, binding port: " + config.WebServiceAddr)
	}
}
