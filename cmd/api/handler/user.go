package handler

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strconv"

	"github.com/bytedance-summer-camp-2023/tiktok/cmd/api/rpc"
	"github.com/bytedance-summer-camp-2023/tiktok/internal/response"
	"github.com/bytedance-summer-camp-2023/tiktok/kitex/kitex_gen/user"
	kitex "github.com/bytedance-summer-camp-2023/tiktok/kitex/kitex_gen/user"
)

type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register 用户注册
func Register(ctx context.Context, c *app.RequestContext) {

	body, err := c.Body()
	if err != nil {
		panic(err)
	}
	var p UserReq
	if err := json.Unmarshal(body, &p); err != nil {

	}
	//校验参数
	if len(p.Username) == 0 || len(p.Password) == 0 {
		c.JSON(http.StatusBadRequest, response.Register{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "用户名或密码不能为空",
			},
		})
		return
	}
	if len(p.Username) > 32 || len(p.Password) > 32 {
		c.JSON(http.StatusOK, response.Register{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "用户名或密码长度不能大于32个字符",
			},
		})
		return
	}
	//调用kitex/kitex_gen
	req := &kitex.UserRegisterRequest{
		Username: p.Username,
		Password: p.Password,
	}
	// 调用rpc user服务进行用户注册
	res, err := rpc.Register(ctx, req)
	if err != nil || res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.Register{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.Register{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		UserID: res.UserId,
		Token:  res.Token,
	})
}

// Login 登录
func Login(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		panic(err)
	}
	var p UserReq
	if err := json.Unmarshal(body, &p); err != nil {

	}
	//校验参数
	if len(p.Username) == 0 || len(p.Password) == 0 {
		c.JSON(http.StatusBadRequest, response.Login{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "用户名或密码不能为空",
			},
		})
		return
	}
	//调用kitex/kitex_gen
	req := &user.UserLoginRequest{
		Username: p.Username,
		Password: p.Password,
	}
	res, _ := rpc.Login(ctx, req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.Login{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, response.Login{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		UserID: res.UserId,
		Token:  res.Token,
	})
}

// UserInfo 用户信息
func UserInfo(ctx context.Context, c *app.RequestContext) {
	userId := c.Query("user_id")
	token := c.Query("token")
	if len(token) == 0 {
		c.JSON(http.StatusOK, response.UserInfo{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "token 已过期",
			},
			User: nil,
		})
		return
	}
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, response.UserInfo{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  "user_id 不合法",
			},
			User: nil,
		})
		return
	}

	//调用kitex/kitex_genit
	req := &kitex.UserInfoRequest{
		UserId: id,
		Token:  token,
	}
	res, _ := rpc.UserInfo(ctx, req)
	if res.StatusCode == -1 {
		c.JSON(http.StatusOK, response.UserInfo{
			Base: response.Base{
				StatusCode: -1,
				StatusMsg:  res.StatusMsg,
			},
			User: nil,
		})
		return
	}
	c.JSON(http.StatusOK, response.UserInfo{
		Base: response.Base{
			StatusCode: 0,
			StatusMsg:  res.StatusMsg,
		},
		User: res.User,
	})
}
