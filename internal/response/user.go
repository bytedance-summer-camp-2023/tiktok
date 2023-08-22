package response

import (
	"github.com/bytedance-summer-camp-2023/tiktok/kitex/kitex_gen/user"
)

type Register struct {
	Base
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type Login struct {
	Base
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfo struct {
	Base
	User *user.User `json:"user"`
}