package main

import (
	"io"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	tb "gopkg.in/tucnak/telebot.v2"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

// NOTIFICATIONTIMEOUT is used for removing system notifications
const NOTIFICATIONTIMEOUT = 10 * time.Second

var db = &gorm.DB{}

func main() {
	logIntoFile := true
	logSQL := true
	var err error
	var loggerFile io.Writer

	if logIntoFile {
		loggerFile, err = os.Create("bot.log")
		check(err)
	} else {
		loggerFile = os.Stdout
	}

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(loggerFile)

	err = godotenv.Load()
	check(err)

	token := os.Getenv("TELEGRAM_TOKEN")

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 30 * time.Second},
	})
	check(err)

	db, err = gorm.Open("sqlite3", "db.sqlite3")
	check(err)
	defer db.Close()

	if logSQL {
		logger := log.StandardLogger()
		logger.Out = loggerFile
		db.SetLogger(logger)
		db.LogMode(true)
	}

	// Migrate the schema
	db.AutoMigrate(&User{}, &LogItem{}, &Category{})

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

	b.Handle("/stats", func(m *tb.Message) {
		_, err := b.Send(m.Sender, "In development 💪")
		check(err)
	})

	b.Handle("/export", func(m *tb.Message) {
		handleExport(m, b)
	})
	b.Handle("/delete", func(m *tb.Message) {
		handleDelete(m, b)
	})

	b.Start()
}
