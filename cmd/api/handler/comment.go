package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"

	"github.com/bytedance-summer-camp-2023/tiktok/cmd/api/rpc"
	"github.com/bytedance-summer-camp-2023/tiktok/internal/response"
	kitex "github.com/bytedance-summer-camp-2023/tiktok/kitex/kitex_gen/comment"
)

func CommentAction(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")

	// query video ID
	vid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.CommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "Invalid video ID",
			},
			Comment: nil,
		})
		return
	}

	// check action type
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		c.JSON(http.StatusOK, response.CommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "Invalid action type",
			},
			Comment: nil,
		})
		return
	}

	// rpc request info
	req := new(kitex.CommentActionRequest)
	req.Token = token
	req.VideoId = vid
	req.ActionType = actionType

	if actionType == 1 {
		// check logic of creating comment
		commentText := c.Query("comment_text")
		if commentText == "" {
			c.JSON(http.StatusOK, response.CommentAction{
				Base: response.Base{
					StatusCode: -1,
					StatusMsg:  "Comment text should not be empty",
				},
				Comment: nil,
			})
			return
		}
		req.CommentText = commentText
	} else if actionType == 2 {
		// check logic of deleting comment
		commentID, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, response.CommentAction{
				Base: response.Base{
					StatusCode: -1,
					StatusMsg:  "Invalid comment ID",
				},
				Comment: nil,
			})
			return
		}
		req.CommentId = commentID
	}

	// rpc request
	res, _ := rpc.CommentAction(ctx, req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.CommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
			Comment: nil,
		})
		return
	}

	// rpc success response
	c.JSON(http.StatusOK, response.CommentAction{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		Comment: res.Comment,
	})
}

func CommentList(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")

	// query video ID
	vid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.CommentAction{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "Invalid video ID",
			},
			Comment: nil,
		})
		return
	}

	// rpc request
	req := &kitex.CommentListRequest{
		Token:   token,
		VideoId: vid,
	}
	res, _ := rpc.CommentList(ctx, req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.CommentList{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
			CommentList: nil,
		})
		return
	}

	// rpc success response
	c.JSON(http.StatusOK, response.CommentList{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		CommentList: res.CommentList,
	})
}
