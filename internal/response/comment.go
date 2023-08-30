package response

import (
	"github.com/bytedance-summer-camp-2023/tiktok/kitex/kitex_gen/comment"
)

type CommentAction struct {
	Base
	Comment *comment.Comment `json:"comment"`
}

type CommentList struct {
	Base
	CommentList []*comment.Comment `json:"comment_list"`
}
