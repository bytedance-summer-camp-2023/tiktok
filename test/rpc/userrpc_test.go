package rpc

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"tiktok/src/constant/config"
	"tiktok/src/rpc/user"
)

func TestGetUserInfo(t *testing.T) {
	var Client user.UserServiceClient
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.UserRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	assert.Empty(t, err)
	Client = user.NewUserServiceClient(conn)
	res, err := Client.GetUserInfo(context.Background(), &user.UserRequest{
		UserId:  16,
		ActorId: 0,
	})
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
}
