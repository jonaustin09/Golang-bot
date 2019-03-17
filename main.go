package main

import (
	"io"
	"os"
	"time"

	"github.com/dobrovolsky/money_bot/stats"

	mb "github.com/dobrovolsky/money_bot/money_bot"
	"google.golang.org/grpc"

	"github.com/jinzhu/gorm"

	tb "gopkg.in/tucnak/telebot.v2"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	var loggerFile io.Writer

	mb.Confg, err = mb.InitConfig()
	mb.Check(err)

	if mb.Confg.LogIntoFile {
		loggerFile, err = os.Create("bot.log")
		mb.Check(err)
	} else {
		loggerFile = os.Stdout
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  mb.Confg.TelegramToken,
		Poller: &tb.LongPoller{Timeout: 30 * time.Second},
	})
	mb.Check(err)

	mb.Db, err = gorm.Open("sqlite3", mb.Confg.DbFile)
	mb.Check(err)
	defer mb.Db.Close()

	if mb.Confg.LogSQL {
		logger := log.StandardLogger()
		logger.Out = loggerFile
		mb.Db.SetLogger(logger)
		mb.Db.LogMode(true)
	}

	grpServerAddress := os.Getenv("GRPC_SERVER_ADDRESS")
	conn, err := grpc.Dial(grpServerAddress, grpc.WithInsecure())
	mb.Check(err)
	defer conn.Close()
	statsClient := stats.NewStatsClient(conn)

	b.Handle("/start", func(m *tb.Message) {
		mb.HandleStart(m, b)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		mb.HandleNewMessage(m, b)
	})

	b.Handle(tb.OnEdited, func(m *tb.Message) {
		mb.HandleEdit(m, b)
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		_, err := b.Send(m.Sender, "Sorry i don't support images 😓")
		mb.Check(err)
	})

	b.Handle("/income", func(m *tb.Message) {
		_, err := b.Send(m.Sender, "In development 💪")
		mb.Check(err)
	})

	b.Handle("/stat_all_by_month", func(m *tb.Message) {
		mb.HandleStatsAllByMonth(m, b, statsClient)
	})

	b.Handle("/stat_all_by_category", func(m *tb.Message) {
		mb.HandleStatsAllByCategory(m, b, statsClient)
	})

	b.Handle("/stat_by_category", func(m *tb.Message) {
		mb.HandleStatsByCategory(m, b, statsClient)
	})

	b.Handle("/export", func(m *tb.Message) {
		mb.HandleExport(m, b)
	})
	b.Handle("/delete", func(m *tb.Message) {
		mb.HandleDelete(m, b)
	})

	b.Start()
}
