package moneybot

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config store configuration params
type Config struct {
	DbFile                     string
	TelegramToken              string
	NotificationTimeout        time.Duration
	MonobankIntegrationEnabled bool
	MonobankWebhookURL         string
	MonobankToken1             string
	MonobankToken2             string
	MonobankPort               int
	MonobankAccount1           string
	MonobankAccount2           string
	UserName1                  string
	UserName2                  string
	SenderID1                  int
	SenderID2                  int
	APIServer                  int
}

// InitConfig init configurations from file and .env
func InitConfig() (*Config, error) {
	v := viper.New()

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	v.SetDefault("db_file", "db.sqlite3")
	v.SetDefault("notification_timeout", 10)
	v.SetDefault("monobank_integration", false)

	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	config.NotificationTimeout = time.Duration(v.GetInt("notification_timeout")) * time.Second
	config.DbFile = v.GetString("db_file")
	config.UserName1 = v.GetString("USERNAME_1")
	config.UserName2 = v.GetString("USERNAME_2")
	config.SenderID1 = v.GetInt("SENDER_ID_1")
	config.SenderID2 = v.GetInt("SENDER_ID_2")
	config.TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	config.APIServer = v.GetInt("API_SERVER_PORT")

	config.MonobankIntegrationEnabled = v.GetBool("monobank_integration")
	if config.MonobankIntegrationEnabled {
		config.MonobankWebhookURL = os.Getenv("MONOBANK_WEBHOOK_URL")
		config.MonobankPort = v.GetInt("MONOBANK_PORT")
		config.MonobankToken1 = os.Getenv("MONOBANK_TOKEN_1")
		config.MonobankToken2 = os.Getenv("MONOBANK_TOKEN_2")
		config.MonobankAccount1 = v.GetString("MONOBANK_ACCOUNT_1")
		config.MonobankAccount2 = v.GetString("MONOBANK_ACCOUNT_2")
	}

	return config, nil
}
