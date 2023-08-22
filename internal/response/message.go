package response

import (
	"github.com/bytedance-summer-camp-2023/tiktok/kitex/kitex_gen/message"
)

type MessageChat struct {
	Base
	MessageList []*message.Message `json:"message_list"`
}

type MessageAction struct {
	Base
}
