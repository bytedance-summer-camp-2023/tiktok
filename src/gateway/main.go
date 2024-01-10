package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"tiktok/src/constant/config"
	"tiktok/src/gateway/about"
)

func main() {
	g := gin.Default()
	g.Use(gzip.Gzip(gzip.DefaultCompression))

	// Register Service
	// Test Service
	g.GET("/about", about.Handle)

	// Production Service
	err := g.Run(config.WebServiceAddr)

	if err != nil {
		panic("Can not run TikTok Gateway, binding port: " + config.WebServiceAddr)
	}
}
