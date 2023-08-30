package main

import (
	"fmt"
	"github.com/bytedance-summer-camp-2023/tiktok/cmd/api/handler"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/viper"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/network/standard"
)

var (
	apiConfig     = viper.Init("api")
	apiServerAddr = fmt.Sprintf("%s:%d", apiConfig.Viper.GetString("server.host"), apiConfig.Viper.GetInt("server.port"))
)

func registerGroup(hz *server.Hertz) {
	douyin := hz.Group("/douyin")
	{
		user := douyin.Group("/user")
		{
			user.GET("/", handler.UserInfo)           //获取用户信息
			user.POST("/register/", handler.Register) //注册接口
			user.POST("/login/", handler.Login)       //登陆接口
		}

		douyin.GET("/feed", handler.Feed)
		publish := douyin.Group("/publish")
		{
			publish.POST("/action/", handler.PublishAction)
			publish.GET("/list/", handler.PublishList)
		}

		favorite := douyin.Group("/favorite")
		{
			favorite.POST("/action/", handler.FavoriteAction)
			favorite.GET("/list/", handler.FavoriteList)
		}

		comment := douyin.Group("/comment")
		{
			comment.POST("/action/", handler.CommentAction)
			comment.GET("/list/", handler.CommentList)
		}

		relation := douyin.Group("/relation")
		{
			// 粉丝列表
			relation.GET("/follower/list/", handler.FollowerList)
			// 关注列表
			relation.GET("/follow/list/", handler.FollowList)
			// 朋友列表
			relation.GET("/friend/list/", handler.FriendList)
			relation.POST("/action/", handler.RelationAction)
		}
	}
}

func InitHertz() *server.Hertz {
	opts := []config.Option{server.WithHostPorts(apiServerAddr)}

	// 配置网络库
	hertzNet := standard.NewTransporter
	opts = append(opts, server.WithTransport(hertzNet))
	hz := server.Default(opts...)

	return hz
}

func main() {
	hz := InitHertz()
	registerGroup(hz)
	hz.Spin()
}
