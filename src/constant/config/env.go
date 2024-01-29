package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

var EnvCfg envConfig

type envConfig struct {
	ConsulAddr            string `env:"CONSUL_ADDR" envDefault:"localhost:8500"`
	ConsulAnonymityPrefix string `env:"CONSUL_ANONYMITY_NAME" envDefault:""`
	LoggerLevel           string `env:"LOGGER_LEVEL" envDefault:"INFO"`
	LoggerWithTraceState  string `env:"LOGGER_OUT_TRACING" envDefault:"disable"`
	TiedLogging           string `env:"TIED" envDefault:"NONE"`
	MySQLHost             string `env:"MYSQL_HOST"`
	MySQLPort             string `env:"MYSQL_PORT"`
	MySQLUser             string `env:"MYSQL_USER"`
	MySQLPassword         string `env:"MYSQL_PASSWORD"`
	MySQLDataBase         string `env:"MYSQL_DATABASE"`
	MySQLSchema           string `env:"MYSQL_SCHEMA" envDefault:""`
	StorageType           string `env:"STORAGE_TYPE" envDefault:"fs"`
	FileSystemStartPath   string `env:"FS_PATH" envDefault:"/tmp"`
	FileSystemBaseUrl     string `env:"FS_BASEURL" envDefault:"http://localhost/"`
	RedisPrefix           string `env:"REDIS_PREFIX" envDefault:""`
	RedisAddr             string `env:"REDIS_ADDR"`
	RedisPassword         string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB               int    `env:"REDIS_DB" envDefault:"0"`
	RedisMaster           string `env:"REDIS_MASTER"`
	TracingEndPoint       string `env:"TRACING_ENDPOINT"`
	PyroscopeState        string `env:"PYROSCOPE_STATE" envDefault:"false"`
	PyroscopeAddr         string `env:"PYROSCOPE_ADDR"`
	RabbitMQUsername      string `env:"RABBITMQ_USERNAME" envDefault:"guest"`
	RabbitMQPassword      string `env:"RABBITMQ_PASSWORD" envDefault:"guest"`
	RabbitMQAddr          string `env:"RABBITMQ_ADDRESS" envDefault:"localhost"`
	RabbitMQPort          string `env:"RABBITMQ_PORT" envDefault:"5672"`
	RabbitMQVhostPrefix   string `env:"RABBITMQ_VHOST_PREFIX" envDefault:""`
	ChatGPTAPIKEYS        string `env:"CHATGPT_API_KEYS"`
}

func init() {
	if err := godotenv.Load("src/constant/config/.env"); err != nil {
		log.Errorf("Can not read env from file system, please check the right this program owned.")
	}

	EnvCfg = envConfig{}

	if err := env.Parse(&EnvCfg); err != nil {
		panic("Can not parse env from file system, please check the env.")
	}
}
