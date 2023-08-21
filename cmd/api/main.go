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
			user.GET("/", handler.UserInfo)
			user.POST("/register/", handler.Register)
			user.POST("/login/", handler.Login)
		}
	}
}

func InitHertz() *server.Hertz {
	//logger := z.InitLogger()
	opts := []config.Option{server.WithHostPorts(apiServerAddr)}
	// 网络库
	hertzNet := standard.NewTransporter
	//if apiConfig.Viper.GetBool("Hertz.useNetPoll") {
	//	hertzNet = netpoll.NewTransporter
	//}
	opts = append(opts, server.WithTransport(hertzNet))
	// TLS & Http2
	// https://github.com/cloudwego/hertz-examples/blob/main/protocol/tls/main.go
	hz := server.Default(opts...)

	return hz
}

func main() {
	hz := InitHertz()
	// add handler
	registerGroup(hz)
	hz.Spin()
}
