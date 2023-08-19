package viper

import (
	V "github.com/spf13/viper"
	"log"
)

type Config struct {
	Viper *V.Viper
}

func Init(configName string) Config {
	config := Config{Viper: V.New()}
	v := config.Viper
	v.SetConfigType("yml")
	v.SetConfigName(configName)
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("errno is %+v", err)
	}
	return config
}
