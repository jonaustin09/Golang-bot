package main

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var db = &gorm.DB{}
var config = &Config{}

type Config struct {
	dbFile              string
	logIntoFile         bool
	logSQL              bool
	telegramToken       string
	GRPCServer          string
	notificationTimeout time.Duration
}

func initConfig() (*Config, error) {
	v := viper.New()

	err := godotenv.Load()
	check(err)

	v.SetDefault("db_file", "db.sqlite3")
	v.SetDefault("enable_file_log", true)
	v.SetDefault("enable_sql_log", true)
	v.SetDefault("notification_timeout", 10)

	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	config.notificationTimeout = time.Duration(v.GetInt("notification_timeout")) * time.Second
	config.logIntoFile = v.GetBool("enable_file_log")
	config.logSQL = v.GetBool("enable_sql_log")
	config.dbFile = v.GetString("db_file")
	config.telegramToken = os.Getenv("TELEGRAM_TOKEN")
	config.GRPCServer = os.Getenv("GRPC_SERVER_ADDRESS")
	config.GRPCServer = os.Getenv("GRPC_SERVER_ADDRESS")

	return config, nil
}
