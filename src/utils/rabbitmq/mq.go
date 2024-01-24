package rabbitmq

import (
	"fmt"
	"tiktok/src/constant/config"
)

func BuildMQConnAddr() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s", config.EnvCfg.RabbitMQUsername, config.EnvCfg.RabbitMQPassword,
		config.EnvCfg.RabbitMQAddr, config.EnvCfg.RabbitMQPassword, config.EnvCfg.RabbitMQVhostPrefix)
}
