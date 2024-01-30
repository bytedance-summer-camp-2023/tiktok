package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net/http"
	"testing"
	"tiktok/src/constant/config"
	"tiktok/src/gateway/models"
	"tiktok/src/rpc/relation"
)

func TestRelationAction(t *testing.T) {

	//url := "http://127.0.0.1:37000/douyin/relation/reg?username=" + uuid.New().String() + "&password=epicmo"
	client := &http.Client{}
	url := "http://127.0.0.1:37000/douyin/relation/action"
	method := "POST"
	req, err := http.NewRequest(method, url, nil)
	q := req.URL.Query()
	q.Add("token", "e65153ea-15b6-4959-9462-f9fb5c5d59ce")
	q.Add("user_id", "1")
	q.Add("actor_id", "4")
	q.Add("action_type", "1")
	req.URL.RawQuery = q.Encode()

	assert.Empty(t, err)

	res, err := client.Do(req)
	assert.Empty(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		assert.Empty(t, err)
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	assert.Empty(t, err)
	relation := &models.RelationActionRes{}
	err = json.Unmarshal(body, &relation)
	assert.Empty(t, err)
	assert.Equal(t, 0, relation.StatusCode)
}

func TestUnfollow(t *testing.T) {
	var Client relation.RelationServiceClient
	req := relation.RelationActionRequest{
		UserId:  2,
		ActorId: 0,
	}

	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.RelationRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	assert.Empty(t, err)
	Client = relation.NewRelationServiceClient(conn)

	res, err := Client.Unfollow(context.Background(), &req)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
}

func TestGetFollowList(t *testing.T) {

	var Client relation.RelationServiceClient
	req := relation.FollowListRequest{
		ActorId: 1,
		UserId:  1,
	}

	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.RelationRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	assert.Empty(t, err)
	Client = relation.NewRelationServiceClient(conn)

	res, err := Client.GetFollowList(context.Background(), &req)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), res.StatusCode)

}

func TestGetFollowerList(t *testing.T) {

	var Client relation.RelationServiceClient
	req := relation.FollowerListRequest{
		ActorId: 1,
		UserId:  1,
	}

	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.RelationRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	assert.Empty(t, err)
	Client = relation.NewRelationServiceClient(conn)

	res, err := Client.GetFollowerList(context.Background(), &req)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
}

func TestCountFollowList(t *testing.T) {
	var Client relation.RelationServiceClient
	req := relation.CountFollowListRequest{
		UserId: 1,
	}
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.RelationRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	assert.Empty(t, err)
	Client = relation.NewRelationServiceClient(conn)

	res, err := Client.CountFollowList(context.Background(), &req)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
}

func TestCountFollowerList(t *testing.T) {

	var Client relation.RelationServiceClient
	req := relation.CountFollowerListRequest{
		UserId: 1,
	}
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.RelationRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	assert.Empty(t, err)
	Client = relation.NewRelationServiceClient(conn)

	res, err := Client.CountFollowerList(context.Background(), &req)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), res.StatusCode)

}

func TestIsFollow(t *testing.T) {

	var Client relation.RelationServiceClient
	req := relation.IsFollowRequest{
		ActorId: 1,
		UserId:  2,
	}
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.RelationRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	assert.Empty(t, err)
	Client = relation.NewRelationServiceClient(conn)

	res, err := Client.IsFollow(context.Background(), &req)
	assert.NoError(t, err)
	assert.Equal(t, true, res.Result)

}

func TestGetFriendList(t *testing.T) {
	var Client relation.RelationServiceClient
	req := relation.FriendListRequest{
		ActorId: 3,
		UserId:  3,
	}
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.RelationRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	assert.Empty(t, err)
	Client = relation.NewRelationServiceClient(conn)

	res, err := Client.GetFriendList(context.Background(), &req)
	fmt.Println(res)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), res.StatusCode)

}
