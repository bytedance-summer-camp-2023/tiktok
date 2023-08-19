package main

import (
	user2 "github.com/bytedance-summer-camp-2023/tiktok/cmd/user/service"
	user "github.com/bytedance-summer-camp-2023/tiktok/kitex_gen/user/userservice"
	"log"
)

func main() {
	svr := user.NewServer(new(user2.UserServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
