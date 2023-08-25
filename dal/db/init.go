// Package db /*
package db

import (
	"fmt"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/zap"
	"gorm.io/gorm/logger"
	"time"

	"github.com/bytedance-summer-camp-2023/tiktok/pkg/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var (
	_db       *gorm.DB
	config    = viper.Init("db")
	zapLogger = zap.InitLogger()
)

func getDsn(driverWithRole string) string {
	username := config.Viper.GetString(fmt.Sprintf("%s.username", driverWithRole))
	password := config.Viper.GetString(fmt.Sprintf("%s.password", driverWithRole))
	host := config.Viper.GetString(fmt.Sprintf("%s.host", driverWithRole))
	port := config.Viper.GetInt(fmt.Sprintf("%s.port", driverWithRole))
	Dbname := config.Viper.GetString(fmt.Sprintf("%s.database", driverWithRole))

	// data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, Dbname)

	return dsn
}

func init() {
	zapLogger.Info("Redis server connection successful!")

	dsn1 := getDsn("mysql.source")

	var err error
	_db, err = gorm.Open(mysql.Open(dsn1), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err.Error())
	}

	// create replicas
	dsn2 := getDsn("mysql.replica1")
	dsn3 := getDsn("mysql.replica2")

	// set dbresolver
	_db.Use(dbresolver.Register(dbresolver.Config{
		Sources:           []gorm.Dialector{mysql.Open(dsn1)},
		Replicas:          []gorm.Dialector{mysql.Open(dsn2), mysql.Open(dsn3)},
		Policy:            dbresolver.RandomPolicy{},
		TraceResolverMode: false,
	}))

	// other tables will be automatically changed
	if err := _db.AutoMigrate(&User{}); err != nil {
		zapLogger.Fatalln(err.Error())
	}

	// create database
	db, err := _db.DB()
	if err != nil {
		zapLogger.Fatalln(err.Error())
	}

	// set database options
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
}

func GetDB() *gorm.DB {
	return _db
}