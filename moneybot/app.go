package moneybot

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

type Application struct {
	Bot               *tb.Bot
	Config            Config
	LogItemRepository LogItemRepository
	IntegrationEvents chan BankData
}

func BuildApp() Application {
	config, err := InitConfig()
	if err != nil {
		panic(err)
	}

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

	db.InstantSet("gorm:auto_preload", true)
	return Application{
		Bot:               b,
		Config:            *config,
		LogItemRepository: NewGormLogItemRepository(db),
		IntegrationEvents: make(chan BankData)}

}

func (a Application) setUpHandlers() {
	a.Bot.Handle("/start", a.handleStart)

	a.Bot.Handle(tb.OnText, a.handleNewMessage)

	a.Bot.Handle(tb.OnEdited, a.handleEdit)

	a.Bot.Handle("/export", a.handleExport)

	a.Bot.Handle("delete", a.handleDelete)
	a.Bot.Handle("d", a.handleDelete)

	go a.handleIntegration(a.IntegrationEvents)

}

func (a Application) setUpIntegrations() {
	if a.Config.MonobankIntegrationEnabled {
		go ListenWebhook(a)

		err := SetWebhook(a.Config.MonobankToken1, a.Config.MonobankWebhookURL, a.Config.MonobankPort)
		if err != nil {
			logrus.Error(err)
		}
		err = SetWebhook(a.Config.MonobankToken2, a.Config.MonobankWebhookURL, a.Config.MonobankPort)
		if err != nil {
			logrus.Error(err)
		}
	}
}

// Start bot and handling requests
func (a Application) Start() {
	go a.startApiServer()

	a.setUpHandlers()
	a.setUpIntegrations()
	a.Bot.Start()
}
