package rpc

import "github.com/bytedance-summer-camp-2023/tiktok/pkg/viper"

func init() {
	// user rpc
	userConfig := viper.Init("user")
	InitUser(&userConfig)

	//video rpc
	videoConfig := viper.Init("video")
	InitVideo(&videoConfig)

	//comment rpc
	commentConfig := viper.Init("comment")
	InitComment(&commentConfig)

	//favorite rpc
	favoriteConfig := viper.Init("favorite")
	InitFavorite(&favoriteConfig)

	// relation rpc
	relationConfig := viper.Init("relation")
	InitRelation(&relationConfig)
}
