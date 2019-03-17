package moneybot

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Db uses setup from main, for making query
var Db = &gorm.DB{}

// Confg uses setup from main
var Confg = &Config{}

// Config store configuration params
type Config struct {
	DbFile              string
	LogIntoFile         bool
	LogSQL              bool
	TelegramToken       string
	GRPCServer          string
	NotificationTimeout time.Duration
}

// InitConfig init configurations from file and .env
func InitConfig() (*Config, error) {
	v := viper.New()

	err := godotenv.Load()
	Check(err)

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
	config.NotificationTimeout = time.Duration(v.GetInt("notification_timeout")) * time.Second
	config.LogIntoFile = v.GetBool("enable_file_log")
	config.LogSQL = v.GetBool("enable_sql_log")
	config.DbFile = v.GetString("db_file")
	config.TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	config.GRPCServer = os.Getenv("GRPC_SERVER_ADDRESS")
	config.GRPCServer = os.Getenv("GRPC_SERVER_ADDRESS")

	return config, nil
}
