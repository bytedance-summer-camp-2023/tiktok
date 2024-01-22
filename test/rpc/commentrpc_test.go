package rpc

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"testing"
	"tiktok/src/constant/config"
	"tiktok/src/rpc/comment"
)

var Client comment.CommentServiceClient

func setup() {
	conn, _ := grpc.Dial(fmt.Sprintf("127.0.0.1%s", config.CommentRpcServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))

	Client = comment.NewCommentServiceClient(conn)
}

func TestActionComment_Add(t *testing.T) {
	res, err := Client.ActionComment(context.Background(), &comment.ActionCommentRequest{
		ActorId:    1,
		VideoId:    0,
		ActionType: comment.ActionCommentType_ACTION_COMMENT_TYPE_ADD,
		Action:     &comment.ActionCommentRequest_CommentText{CommentText: "Test comment"},
	})
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
}

func TestActionComment_Delete(t *testing.T) {
	res, err := Client.ActionComment(context.Background(), &comment.ActionCommentRequest{
		ActorId:    1,
		VideoId:    0,
		ActionType: comment.ActionCommentType_ACTION_COMMENT_TYPE_DELETE,
		Action:     &comment.ActionCommentRequest_CommentId{CommentId: 1},
	})
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
}

func TestListComment(t *testing.T) {
	res, err := Client.ListComment(context.Background(), &comment.ListCommentRequest{
		ActorId: 1,
		VideoId: 0,
	})
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
}

func TestCountComment(t *testing.T) {
	res, err := Client.CountComment(context.Background(), &comment.CountCommentRequest{
		ActorId: 1,
		VideoId: 0,
	})
	assert.Empty(t, err)
	assert.Equal(t, int32(0), res.StatusCode)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}
