package service

import (
	"tiktok/internal/tool"
	"tiktok/pkg/jwt"
	"tiktok/pkg/rabbitmq"
	"tiktok/pkg/viper"
	"tiktok/pkg/zap"
)

var (
	Jwt        *jwt.JWT
	logger     = zap.InitLogger()
	config     = viper.Init("rabbitmq")
	autoAck    = config.Viper.GetBool("consumer.relation.autoAck")
	RelationMq = rabbitmq.NewRabbitMQSimple("relation", autoAck)
	err        error
	privateKey string
)

func Init(signingKey string) {
	Jwt = jwt.NewJWT([]byte(signingKey))
	privateKey, _ = tool.ReadKeyFromFile(tool.PrivateKeyFilePath)

	go consume()
}
