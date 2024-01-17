package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"tiktok/src/constant/config"
	"tiktok/src/models"
	"tiktok/src/utils/logging"
)

var Client *gorm.DB

func init() {
	var err error
	gormLogrus := logging.GetGormLogger()
	// data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.EnvCfg.MySQLUser,
		config.EnvCfg.MySQLPassword,
		config.EnvCfg.MySQLHost,
		config.EnvCfg.MySQLPort,
		config.EnvCfg.MySQLDataBase)
	if Client, err = gorm.Open(
		mysql.Open((dsn)), &gorm.Config{
			PrepareStmt: true,
			Logger:      gormLogrus,
		},
	); err != nil {
		panic(err)
	}

	if err := Client.AutoMigrate(&models.User{}); err != nil {
		panic(err)
	}

	if err := Client.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
}
