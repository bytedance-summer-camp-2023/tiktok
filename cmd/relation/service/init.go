package service

import (
	"github.com/bytedance-summer-camp-2023/tiktok/internal/tool"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/jwt"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/rabbitmq"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/viper"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/zap"
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
