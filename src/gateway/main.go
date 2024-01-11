package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"tiktok/src/constant/config"
	"tiktok/src/gateway/about"
	"tiktok/src/gateway/auth"
	"tiktok/src/gateway/authmw"
	"tiktok/src/gateway/middleware"
)

func main() {
	g := gin.Default()
	// Configure Gzip
	g.Use(gzip.Gzip(gzip.DefaultCompression))
	// Configure Tracing
	g.Use(middleware.Jaeger())
	g.Use(authmw.TokenAuthMiddleware())
	// Register Service
	// Test Service
	g.GET("/about", about.Handle)

	// Production Service
	rootPath := g.Group("/douyin")
	user := rootPath.Group("/user")
	{
		user.GET("/login", auth.LoginHandle)
		user.GET("/register", auth.RegisterHandle)
	}

	// Run Server
	err := g.Run(config.WebServiceAddr)

	if err != nil {
		panic("Can not run GuGoTik Gateway, binding port: " + config.WebServiceAddr)
	}
}
