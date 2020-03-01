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
	statsClient := stats.NewStatsClient(conn)

	logItemRepository := mb.NewGormLogItemRepository(db)

	b.Handle("/start", func(m *tb.Message) {
		mb.HandleStart(m, b, config)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		mb.HandleNewMessage(m, b, logItemRepository, config)
	})

	b.Handle(tb.OnEdited, func(m *tb.Message) {
		mb.HandleEdit(m, b, logItemRepository, config)
	})

	b.Handle("/stat_all_by_month", func(m *tb.Message) {
		mb.HandleStatsAllByMonth(m, b, statsClient, logItemRepository, config)
	})

	b.Handle("/stat_current_month", func(m *tb.Message) {
		mb.HandleStatsByCategoryForCurrentMonth(m, b, statsClient, logItemRepository, config)
	})

	b.Handle("/export", func(m *tb.Message) {
		mb.HandleExport(m, b, logItemRepository, config)
	})
	b.Handle("delete", func(m *tb.Message) {
		mb.HandleDelete(m, b, logItemRepository, config)
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		err := mb.SendDeletableMessage(m.Sender, b, "Sorry i don't support images 😓", config.NotificationTimeout)
		if err != nil {
			log.Error(err)
		}
	})

	b.Handle("/income", func(m *tb.Message) {
		err := mb.SendDeletableMessage(m.Sender, b, "In development 💪", config.NotificationTimeout)
		if err != nil {
			log.Error(err)
		}
	})

	integrationEvents := make(chan mb.Item)
	go mb.HandleIntegration(integrationEvents, b, logItemRepository, config)

	if config.MonobankIntegrationEnabled {
		go mb.ListenWebhook(config.MonobankPort, integrationEvents)
		err := mb.SetWebhook(config.MonobankToken, config.MonobankWebhookURL, config.MonobankPort)
		if err != nil {
			log.Error(err)
		}
	}

	b.Start()
}
