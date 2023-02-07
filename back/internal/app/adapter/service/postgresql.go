package service

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/spf13/viper"
)

const (
	MaxIdleConnections = 2
	MaxOpenConnections = 10
)

func GetPostgresqlDB() (db *gorm.DB) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		viper.Get(KeyPostgresqlHost),
		viper.Get(KeyPostgresqlPort),
		viper.Get(KeyPostgresqlUser),
		viper.Get(KeyPostgresqlPassword),
		viper.Get(KeyPostgresqlDBName),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(MaxIdleConnections)
	sqlDB.SetMaxOpenConns(MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
