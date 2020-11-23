package moneybot

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dobrovolsky/money_bot/stats"

	"github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Application struct {
	Bot               *tb.Bot
	Config            Config
	LogItemRepository LogItemRepository
	StatsClient       stats.StatsClient
	IntegrationEvents chan Item
}

func (a Application) setUpHandlers() {
	a.Bot.Handle("/start", a.handleStart)

	a.Bot.Handle(tb.OnText, a.handleNewMessage)

	a.Bot.Handle(tb.OnEdited, a.handleEdit)

	a.Bot.Handle("/stat_all_time", a.handleStatAsTableForAllTime)

	a.Bot.Handle("/stat_last_year", a.handleStatAsTableForPeriod)

	a.Bot.Handle("/stat_current_month", a.handleStatsGroupByCategoryForCurrentMonth)

	a.Bot.Handle("/export", a.handleExport)

	a.Bot.Handle("delete", a.handleDelete)
	a.Bot.Handle("d", a.handleDelete)

	go a.handleIntegration(a.IntegrationEvents)

}

func (a Application) setUpIntegrations() {
	if a.Config.MonobankIntegrationEnabled {
		go ListenWebhook(a.Config.MonobankPort, a.IntegrationEvents)

		err := SetWebhook(a.Config.MonobankToken, a.Config.MonobankWebhookURL, a.Config.MonobankPort)
		if err != nil {
			logrus.Error(err)
		}
	}

}

// Start handles requests
func (a Application) Start() {
	go a.startApiServer()

	a.setUpHandlers()
	a.setUpIntegrations()
	a.Bot.Start()
}

// handleStart greeting, saves information about user
func (a Application) handleStart(m *tb.Message) {
	logrus.Infof("Start handleStart request with %s by %v", m.Text, m.Sender.ID)

	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	text := "Hello there i'll help you with your finances! \n" +
		"Use the following format: `item amount`. *For example*: tea 10 (repository name) \n" +
		"To delete message start to replay what you want to delete and type button 'delete'"

	err := SendDeletableMessage(m.Sender, a.Bot, text, a.Config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
		return
	}
}

// handleNewMessage process new messages
func (a Application) handleNewMessage(m *tb.Message) {
	logrus.Infof("Start handleNewMessage request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	go Notify(m.Sender, a.Bot, tb.Typing)

	items := GetItem(m.Text)
	logrus.Info("Parsed data", items)

	if len(items) == 0 {
		text := "Use the following format: `item amount`. *For example*: tea 10 (category name)"
		err := SendDeletableMessage(m.Sender, a.Bot, text, a.Config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}

		go DeleteMessage(m, a.Bot, a.Config.NotificationTimeout)

		return

	}

	if m.ReplyTo != nil {
		if !a.LogItemRepository.RecordExists(int32(m.ReplyTo.ID)) {
			text := "You can not edit this message"
			err := SendDeletableMessage(m.Sender, a.Bot, text, a.Config.NotificationTimeout)
			if err != nil {
				logrus.Error(err)
				return
			}
		} else {
			logrus.Info("Start editing")

			go DeleteMessage(m.ReplyTo, a.Bot, 0)
			err := a.LogItemRepository.DeleteRecordsByMessageID(int32(m.ReplyTo.ID))
			if err != nil {
				logrus.Error(err)
			}

			logrus.Info("Remove all related records")

			SaveItems(items, int32(m.ID), m.Sender, a.Bot, a.LogItemRepository, a.Config)
		}

	} else {
		SaveItems(items, int32(m.ID), m.Sender, a.Bot, a.LogItemRepository, a.Config)
	}
}

// handleEdit allow to edit infromation from db for following message
func (a Application) handleEdit(m *tb.Message) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	go Notify(m.Sender, a.Bot, tb.Typing)

	item := GetItem(m.Text)
	logrus.Info("Parsed data", item)

	a.editLogs(int32(m.ID), m.Sender, item)
}

// editLogs tries to edit log
func (a Application) editLogs(messageID int32, sender *tb.User, items []Item) {
	logrus.Info("Start editing")
	var text string
	var err error

	if len(items) == 0 {
		text = "Use the following format: `item amount`. *For example*: tea 10 (category name)"
		err = SendDeletableMessage(sender, a.Bot, text, a.Config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return

	}

	err = a.LogItemRepository.DeleteRecordsByMessageID(messageID)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("Remove all related records")

	SaveItems(items, messageID, sender, a.Bot, a.LogItemRepository, a.Config)
}

// handleDelete allow to delete infromation from db for following message
func (a Application) handleDelete(m *tb.Message) {
	logrus.Infof("Start handleDelete request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	if m.ReplyTo == nil {
		text := "You should reply for a message which you want to delete ↩️"
		err := SendDeletableMessage(m.Sender, a.Bot, text, a.Config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
			return
		}
	} else {
		err := a.LogItemRepository.DeleteRecordsByMessageID(int32(m.ReplyTo.ID))
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Info("Remove all related records")

		text := "`Remove item`"
		err = SendDeletableMessage(m.Sender, a.Bot, text, a.Config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
			return
		}

		go DeleteMessage(m, a.Bot, a.Config.NotificationTimeout)
		if m.ReplyTo != nil {
			go DeleteMessage(m.ReplyTo, a.Bot, a.Config.NotificationTimeout)
		}

	}
}

// handleStatAsTableForAllTime allow to get information grouped by month and categories as table all stats
func (a Application) handleStatAsTableForAllTime(m *tb.Message) {
	logrus.Infof("Start handleStatsAllByMonth request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	go Notify(m.Sender, a.Bot, tb.UploadingDocument)

	items, err := a.LogItemRepository.GetAggregatedRecordsAllTime()
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		err = SendDeletableMessage(m.Sender, a.Bot, "There are not any records yet 😒", a.Config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	itemsForAnalyze := PrepareForAnalyze(items)

	logrus.Info("Call GetMonthAmountStat")
	monthAmountStat, err := a.StatsClient.GetStatAsTable(context.Background(), &stats.LogItemQueryMessage{
		LogMessagesAggregated: itemsForAnalyze,
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	fileName := fmt.Sprintf("%v-%v-stat.pdf", m.Sender.ID, Timestamp())
	err = SendDocumentFromReader(m.Sender, a.Bot, fileName, monthAmountStat.Res, a.Config)
	if err != nil {
		logrus.Error(err)
	}

	go DeleteMessage(m, a.Bot, a.Config.NotificationTimeout)
}

// handleStatAsTableForPeriod allow to get information grouped by month and categories as table for period of time
func (a Application) handleStatAsTableForPeriod(m *tb.Message) {
	logrus.Infof("Start handleStatAsTableForPeriod request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	go Notify(m.Sender, a.Bot, tb.UploadingPhoto)

	now := time.Now()
	startFromGivenYear, startFromGivenMonth, _ := now.AddDate(-1, 0, 0).Date()
	currentLocation := now.Location()
	firstDayOfMonth := time.Date(startFromGivenYear, startFromGivenMonth, 1, 0, 0, 0, 0, currentLocation)

	items, err := a.LogItemRepository.GetAggregatedRecordsPeriod(firstDayOfMonth.Unix(), now.Unix())

	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		err = SendDeletableMessage(m.Sender, a.Bot, "There are not any records yet 😒", a.Config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	itemsForAnalyze := PrepareForAnalyze(items)

	logrus.Info("Call GetMonthAmountStat")
	monthAmountStat, err := a.StatsClient.GetStatAsTable(context.Background(), &stats.LogItemQueryMessage{
		LogMessagesAggregated: itemsForAnalyze,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	fileName := fmt.Sprintf("%v-%v-stat.pdf", m.Sender.ID, Timestamp())
	err = SendDocumentFromReader(m.Sender, a.Bot, fileName, monthAmountStat.Res, a.Config)
	if err != nil {
		logrus.Error(err)
	}

	go DeleteMessage(m, a.Bot, a.Config.NotificationTimeout)
}

// handleStatsGroupByCategoryForCurrentMonth allow to get information grouped by categories for current month
func (a Application) handleStatsGroupByCategoryForCurrentMonth(m *tb.Message) {
	logrus.Infof("Start HandleStatsByCategoryForCurrentMonth request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	go Notify(m.Sender, a.Bot, tb.UploadingPhoto)

	items, err := a.LogItemRepository.GetAggregatedRecordsCurrentMonth()
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		err = SendDeletableMessage(m.Sender, a.Bot, "There are not any records yet 😒", a.Config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	logrus.Info("Call GetCategoryStat")
	stat, err := a.StatsClient.GetStatGroupByCategory(context.Background(), &stats.LogItemQueryMessage{
		LogMessagesAggregated: PrepareForAnalyze(items),
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	photo := &tb.Photo{File: tb.FromReader(bytes.NewReader(stat.Res))}

	err = SendDeletableMessage(m.Sender, a.Bot, photo, a.Config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
	}

	go DeleteMessage(m, a.Bot, a.Config.NotificationTimeout)
}

// handleExport allow to export data into csv file
func (a Application) handleExport(m *tb.Message) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	go Notify(m.Sender, a.Bot, tb.UploadingDocument)

	items, err := a.LogItemRepository.GetRecords()
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		err = SendDeletableMessage(m.Sender, a.Bot, "There are not any records yet 😒", a.Config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	fileName := fmt.Sprintf("%v-%v-export.csv", m.Sender.ID, Timestamp())
	file, err := os.Create(fileName)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("Create file")

	defer func() {
		err := os.Remove(fileName)
		if err != nil {
			logrus.Error(err)
		}
	}()
	defer func() {
		err := file.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	writer := csv.NewWriter(file)

	for _, item := range items {
		err = writer.Write(item.toCSV())
		if err != nil {
			logrus.Error(err)
			return
		}
	}
	writer.Flush()
	logrus.Info("Save file")

	document := &tb.Document{File: tb.FromDisk(fileName)}

	err = SendDeletableMessage(m.Sender, a.Bot, document, a.Config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Send file to ", m.Sender.ID)

	go DeleteMessage(m, a.Bot, a.Config.NotificationTimeout)
}

// handleIntegration allows to add new integration for example bank
func (a Application) handleIntegration(items <-chan Item) {
	recipient := User{
		ID: a.Config.ChatID,
	}
	var err error
	for item := range items {
		if item.IsValid() {
			if item.Category == "" {
				item.Category, err = a.LogItemRepository.FetchMostRelevantCategory(item.Name)
				if err != nil {
					logrus.Error(err)
				}
			}

			text := fmt.Sprintf("%s %.2f %s", item.Name, item.Amount, item.Category)
			message, err := SendMessage(recipient, a.Bot, text)

			if err != nil {
				logrus.Error(err)
				continue
			}

			_, err = item.ProcessSaving(int32(message.ID), recipient.ID, a.Bot, a.LogItemRepository, a.Config)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}

// startApiServer starts simple http server to return export
func (a Application) startApiServer(){
	http.HandleFunc("/api/export", func(w http.ResponseWriter, r *http.Request) {

		items, err := a.LogItemRepository.GetRecords()
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			return
		}

		logrus.Infof("Fetch items count %v", len(items))

		serializedItems, err := json.Marshal(items)
		w.Header().Set("Content-Type", "application/json")
		w.Write(serializedItems)
	})

	logrus.Info(fmt.Sprintf("listen api on: %d", a.Config.APIServer))
	err := http.ListenAndServe(fmt.Sprintf(":%d", a.Config.APIServer), nil)
	if err != nil {
		logrus.Error(err)
	}
}
