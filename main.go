package main

import (
	"io"
	"os"
	"time"

	"github.com/dobrovolsky/money_bot/stats"

	mb "github.com/dobrovolsky/money_bot/moneybot"
	"google.golang.org/grpc"

	"github.com/jinzhu/gorm"

	tb "gopkg.in/tucnak/telebot.v2"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

func main() {

	var err error
	var loggerFile io.Writer

	_config, err := mb.InitConfig()
	if err != nil {
		panic(err)
	}

	config := *_config

	if config.LogIntoFile {
		loggerFile, err = os.Create("bot.log")
		if err != nil {
			log.Error(err)
		}
	} else {
		loggerFile = os.Stdout
	}

	logger := log.StandardLogger()
	logger.Out = loggerFile

	b, err := tb.NewBot(tb.Settings{
		Token:  config.TelegramToken,
		Poller: &tb.LongPoller{Timeout: 30 * time.Second},
	})
	if err != nil {
		log.Error(err)
	}

	db, err := gorm.Open("sqlite3", config.DbFile)
	if err != nil {
		log.Error(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	db.InstantSet("gorm:auto_preload", true)

	if config.LogSQL {
		db.SetLogger(logger)
		db.LogMode(true)
	}

	grpServerAddress := os.Getenv("GRPC_SERVER_ADDRESS")
	conn, err := grpc.Dial(grpServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Error(err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	app := mb.Application{
		Bot:               b,
		Config:            config,
		LogItemRepository: mb.NewGormLogItemRepository(db),
		StatsClient:       stats.NewStatsClient(conn),
		IntegrationEvents: make(chan mb.Item)}

	app.Start()
}
