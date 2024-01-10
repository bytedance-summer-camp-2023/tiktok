package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

var EnvCfg envConfig

type envConfig struct {
	ConsulAddr          string `env:"CONSUL_ADDR" envDefault:"localhost:8500"`
	LoggerLevel         string `env:"LOGGER_LEVEL" envDefault:"INFO"`
	TiedLogging         string `env:"TIED" envDefault:"NONE"`
	MySQLHost           string `env:"MySQL_HOST"`
	MySQLPort           string `env:"MySQL_PORT"`
	MySQLUser           string `env:"MySQL_USER"`
	MySQLPassword       string `env:"MySQL_PASSWORD"`
	MySQLDataBase       string `env:"MySQL_DATABASE"`
	StorageType         string `env:"STORAGE_TYPE" envDefault:"fs"`
	FileSystemStartPath string `env:"FS_PATH" envDefault:"/tmp"`
	FileSystemBaseUrl   string `env:"FS_BASEURL" envDefault:"http://localhost/"`
	RedisAddr           string `env:"REDIS_ADDR"`
	RedisPassword       string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB             int    `env:"REDIS_DB" envDefault:"0"`
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
