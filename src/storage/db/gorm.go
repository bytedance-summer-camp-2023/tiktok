package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"gorm.io/plugin/opentelemetry/tracing"
	"strings"
	"tiktok/src/constant/config"
	"tiktok/src/utils/logging"
	"time"
)

var Client *gorm.DB

func init() {
	var err error

	gormLogrus := logging.GetGormLogger()

	var cfg gorm.Config
	if config.EnvCfg.MySQLSchema == "" {
		cfg = gorm.Config{
			PrepareStmt: true,
			Logger:      gormLogrus,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: config.EnvCfg.MySQLSchema + "." + config.EnvCfg.MySQLPrefix,
			},
		}
	} else {
		cfg = gorm.Config{
			PrepareStmt: true,
			Logger:      gormLogrus,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: config.EnvCfg.MySQLSchema + "." + config.EnvCfg.MySQLPrefix,
			},
		}
	}

	if Client, err = gorm.Open(
		mysql.Open(
			fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
				config.EnvCfg.MySQLHost,
				config.EnvCfg.MySQLUser,
				config.EnvCfg.MySQLPassword,
				config.EnvCfg.MySQLDataBase,
				config.EnvCfg.MySQLPort)),
		&cfg,
	); err != nil {
		panic(err)
	}

	if config.EnvCfg.MySQLReplicaState == "enable" {
		var replicas []gorm.Dialector
		for _, addr := range strings.Split(config.EnvCfg.MySQLReplicaAddress, ",") {
			pair := strings.Split(addr, ":")
			if len(pair) != 2 {
				continue
			}

			replicas = append(replicas, mysql.Open(
				fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
					pair[0],
					config.EnvCfg.MySQLReplicaUsername,
					config.EnvCfg.MySQLReplicaPassword,
					config.EnvCfg.MySQLDataBase,
					pair[1])))
		}

		err := Client.Use(dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}))
		if err != nil {
			panic(err)
		}
	}

	sqlDB, err := Client.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(24 * time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Hour)

	if err := Client.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
}
