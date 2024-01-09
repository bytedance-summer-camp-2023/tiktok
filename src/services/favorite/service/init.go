package service

import (
	"tiktok/utils/jwt"
	"tiktok/utils/rabbitmq"
	"tiktok/utils/viper"
	"tiktok/utils/zap"
)

var (
	Jwt        *jwt.JWT
	logger     = zap.InitLogger()
	config     = viper.Init("rabbitmq")
	autoAck    = config.Viper.GetBool("consumer.favorite.autoAck")
	FavoriteMq = rabbitmq.NewRabbitMQSimple("favorite", autoAck)
	err        error
)

func Init(signingKey string) {
	Jwt = jwt.NewJWT([]byte(signingKey))
	//GoCron()
	go consume()
}
