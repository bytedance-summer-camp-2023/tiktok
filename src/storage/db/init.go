package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
	"tiktok/src/constant/config"
	"tiktok/src/utils/logging"
	"time"
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
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: config.EnvCfg.MySQLSchema,
			},
		},
	); err != nil {
		panic(err)
	}

	sqlDB, err := Client.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := Client.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
}
