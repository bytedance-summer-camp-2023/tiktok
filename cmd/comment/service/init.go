package service

import (
	"tiktok/pkg/jwt"
)

var (
	Jwt *jwt.JWT
)

func Init(signingKey string) {
	Jwt = jwt.NewJWT([]byte(signingKey))
}
