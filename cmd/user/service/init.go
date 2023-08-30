package service

import (
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/jwt"
)

var (
	Jwt *jwt.JWT
)

func Init(signingKey string) {
	Jwt = jwt.NewJWT([]byte(signingKey))
}
