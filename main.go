package main

import (
	"io"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/spf13/viper"

	"github.com/dobrovolsky/money_bot/stats"
	"google.golang.org/grpc"

	"github.com/jinzhu/gorm"

	tb "gopkg.in/tucnak/telebot.v2"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

// NOTIFICATIONTIMEOUT is used for removing system notifications
const NOTIFICATIONTIMEOUT = 10 * time.Second

var db = &gorm.DB{}
var config = &Config{}

type Config struct {
	dbFile        string
	logIntoFile   bool
	logSQL        bool
	telegramToken string
	GRPCServer    string
}

func initConfig() (*Config, error) {
	v := viper.New()

	err := godotenv.Load()
	check(err)

	v.SetDefault("db_file", "db.sqlite3")
	v.SetDefault("enable_file_log", true)
	v.SetDefault("enable_sql_log", true)

	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	config.logIntoFile = v.GetBool("enable_file_log")
	config.logSQL = v.GetBool("enable_sql_log")
	config.dbFile = v.GetString("db_file")
	config.telegramToken = os.Getenv("TELEGRAM_TOKEN")
	config.GRPCServer = os.Getenv("GRPC_SERVER_ADDRESS")

	return config, nil
}

func main() {
	var err error
	var loggerFile io.Writer

	config, err = initConfig()
	check(err)

	if config.logIntoFile {
		loggerFile, err = os.Create("bot.log")
		check(err)
	} else {
		loggerFile = os.Stdout
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  config.telegramToken,
		Poller: &tb.LongPoller{Timeout: 30 * time.Second},
	})
	check(err)

	db, err = gorm.Open("sqlite3", config.dbFile)
	check(err)
	defer db.Close()

	if config.logSQL {
		logger := log.StandardLogger()
		logger.Out = loggerFile
		db.SetLogger(logger)
		db.LogMode(true)
	}

	// Migrate the schema
	db.AutoMigrate(&User{}, &LogItem{}, &Category{})

	grpServerAddress := os.Getenv("GRPC_SERVER_ADDRESS")
	conn, err := grpc.Dial(grpServerAddress, grpc.WithInsecure())
	check(err)
	defer conn.Close()
	statsClient := stats.NewStatsClient(conn)

	b.Handle("/start", func(m *tb.Message) {
		handleStart(m, b)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		handleNewMessage(m, b)
	})

	b.Handle(tb.OnEdited, func(m *tb.Message) {
		handleEdit(m, b)
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		_, err := b.Send(m.Sender, "Sorry i don't support images 😓")
		check(err)
	})

	b.Handle("/income", func(m *tb.Message) {
		_, err := b.Send(m.Sender, "In development 💪")
		check(err)
	})

	b.Handle("/stat_all_by_month", func(m *tb.Message) {
		handleStatsAllByMonth(m, b, statsClient)

	})

	b.Handle("/stat_all_by_category", func(m *tb.Message) {
		handleStatsAllByCategory(m, b, statsClient)
	})

	b.Handle("/export", func(m *tb.Message) {
		handleExport(m, b)
	})
	b.Handle("/delete", func(m *tb.Message) {
		handleDelete(m, b)
	})

	b.Start()
}
