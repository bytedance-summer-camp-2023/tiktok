package rpc

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"tiktok/kitex/kitex_gen/comment"
	"tiktok/kitex/kitex_gen/comment/commentservice"
	"tiktok/pkg/etcd"
	"tiktok/pkg/middleware"
	"tiktok/pkg/viper"
	"time"
)

var (
	commentClient commentservice.Client
)

func InitComment(config *viper.Config) {
	etcdAddr := fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	serviceName := config.Viper.GetString("server.name")
	r, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		panic(err)
	}

	c, err := commentservice.NewClient(
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
	commentClient = c
}

func CommentAction(ctx context.Context, req *comment.CommentActionRequest) (*comment.CommentActionResponse, error) {
	return commentClient.CommentAction(ctx, req)
}

func CommentList(ctx context.Context, req *comment.CommentListRequest) (*comment.CommentListResponse, error) {
	return commentClient.CommentList(ctx, req)
}
