package main

import (
	"fmt"
	"github.com/bytedance-summer-camp-2023/tiktok/cmd/video/service"
	video "github.com/bytedance-summer-camp-2023/tiktok/kitex_gen/video/videoservice"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/etcd"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/middleware"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/viper"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/zap"
	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

var (
	config      = viper.Init("video")
	serviceName = config.Viper.GetString("server.name")
	serviceAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
	etcdAddr    = fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	signingKey  = config.Viper.GetString("JWT.signingKey")
	logger      = zap.InitLogger()
)

func init() {
	service.Init(signingKey)
}

func main() {

	// 服务注册
	r, err := etcd.NewEtcdRegistry([]string{etcdAddr})
	if err != nil {
		logger.Fatalln(err.Error())
	}

	addr, err := net.ResolveTCPAddr("tcp", serviceAddr)
	if err != nil {
		logger.Fatalln(err.Error())
	}

	svr := video.NewServer(new(service.VideoServiceImpl),
		server.WithServiceAddr(addr),
		server.WithMiddleware(middleware.CommonMiddleware),
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithRegistry(r),
		//server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		// server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)

	if err := svr.Run(); err != nil {
		logger.Fatalf("%v stopped with error: %v", serviceName, err.Error())
	}
}
