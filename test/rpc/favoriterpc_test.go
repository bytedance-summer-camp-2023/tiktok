package rpc

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"tiktok/src/constant/config"
	"tiktok/src/rpc/favorite"
)

var favoriteClient favorite.FavoriteServiceClient

func setups1() {
	conn, _ := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.FavoriteRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	favoriteClient = favorite.NewFavoriteServiceClient(conn)
}
func TestFavoriteAction(t *testing.T) {
	setups1()
	res, err := favoriteClient.FavoriteAction(context.Background(), &favorite.FavoriteRequest{
		ActorId:    2,
		VideoId:    1,
		ActionType: 1,
	})
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
}

func TestFavoriteList(t *testing.T) {
	setups1()
	res, err := favoriteClient.FavoriteList(context.Background(), &favorite.FavoriteListRequest{
		ActorId: 1,
		UserId:  1,
	})

	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Nil(t, res.VideoList)
}

func TestIsFavorite(t *testing.T) {
	setups1()
	res, err := favoriteClient.IsFavorite(context.Background(), &favorite.IsFavoriteRequest{
		ActorId: 1,
		VideoId: 1,
	})
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, true, res.Result)
}

func TestCountFavorite(t *testing.T) {
	setups1()
	res, err := favoriteClient.CountFavorite(context.Background(), &favorite.CountFavoriteRequest{
		VideoId: 1,
	})
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, uint32(1), res.Count)
}

func TestCountUserFavorite(t *testing.T) {
	setups1()
	res, err := favoriteClient.CountUserFavorite(context.Background(), &favorite.CountUserFavoriteRequest{
		UserId: 1,
	})
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, uint32(1), res.Count)
}

func TestCountUserTotalFavorited(t *testing.T) {
	setups1()
	res, err := favoriteClient.CountUserTotalFavorited(context.Background(), &favorite.CountUserTotalFavoritedRequest{
		ActorId: 1,
		UserId:  1,
	})
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, uint32(1), res.Count)
}
