package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	video "tiktok/kitex/kitex_gen/video"
	"tiktok/kitex/kitex_gen/video/videoservice"
	"tiktok/pkg/etcd"
	"tiktok/pkg/middleware"
	"tiktok/pkg/viper"
)

var (
	videoClient videoservice.Client
)

func InitVideo(config *viper.Config) {
	etcdAddr := fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	serviceName := config.Viper.GetString("server.name")
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}

	c, err := videoservice.NewClient(
		serviceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(1),                        // mux
		client.WithRPCTimeout(300*time.Second),             // rpc timeout
		client.WithConnectTimeout(300000*time.Millisecond), // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()),  // retry
		client.WithResolver(r),                             // resolver
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)
	if err != nil {
		panic(err)
	}
	videoClient = c
}

func Feed(ctx context.Context, req *video.FeedRequest) (*video.FeedResponse, error) {
	return videoClient.Feed(ctx, req)
}

func PublishAction(ctx context.Context, req *video.PublishActionRequest) (*video.PublishActionResponse, error) {
	return videoClient.PublishAction(ctx, req)
}

func PublishList(ctx context.Context, req *video.PublishListRequest) (*video.PublishListResponse, error) {
	return videoClient.PublishList(ctx, req)
}
