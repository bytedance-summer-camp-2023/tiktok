package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
	"tiktok/src/constant/config"
	"tiktok/src/models"
	"time"
)

var DB *gorm.DB
var Client *gorm.DB

func init() {
	var err error
	gormLogrus := logger.New(
		logrus.NewWriter(),
		logger.Config{
			SlowThreshold: time.Millisecond,
			Colorful:      false,
			LogLevel:      logger.Info,
		},
	)

	if DB, err = gorm.Open(
		mysql.Open(
			fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
				config.EnvCfg.MySQLHost,
				config.EnvCfg.MySQLUser,
				config.EnvCfg.MySQLPassword,
				config.EnvCfg.MySQLDataBase,
				config.EnvCfg.MySQLPort)),
		&gorm.Config{
			PrepareStmt: true,
			Logger:      gormLogrus,
		},
	); err != nil {
		panic(err)
	}

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		panic(err)
	}

	if err := DB.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
}
