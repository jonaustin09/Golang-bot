package main

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	tb "gopkg.in/tucnak/telebot.v2"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db = &gorm.DB{}

func main() {
	err := godotenv.Load()
	Check(err)

	token := os.Getenv("TELEGRAM_TOKEN")

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 3 * time.Second},
	})
	Check(err)

	db, err = gorm.Open("sqlite3", "db.sqlite3")
	Check(err)
	defer db.Close()

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
		Check(err)
	})

	b.Handle("/income", func(m *tb.Message) {
		_, err := b.Send(m.Sender, "In development 💪")
		Check(err)
	})
	b.Handle("/export", func(m *tb.Message) {
		_, err := b.Send(m.Sender, "In development 💪")
		Check(err)
		handleExport(m, b)
	})

	b.Start()
}
