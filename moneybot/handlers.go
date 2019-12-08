package moneybot

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/dobrovolsky/money_bot/stats"

	"github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

// HandleStart greeting, saves information about user
func HandleStart(m *tb.Message, b *tb.Bot, config Config) {
	logrus.Infof("Start handleStart request with %s by %v", m.Text, m.Sender.ID)

	if isForbidden(m, b, config) {
		return
	}

	text := "Hello there i'll help you with your finances! \n" +
		"Use the following format: `item amount`. *For example*: tea 10 (repository name) \n" +
		"To delete message start to replay what you want to delete and type button 'delete'"

	err := SendDeletableMessage(m.Sender, b, text, config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
		return
	}
}

// HandleNewMessage process new messages
func HandleNewMessage(m *tb.Message, b *tb.Bot, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleNewMessage request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, b, config) {
		return
	}

	go Notify(m.Sender, b, tb.Typing)

	items := GetItem(m.Text)
	logrus.Info("Parsed data", items)

	if len(items) == 0 {
		text := "Use the following format: `item amount`. *For example*: tea 10 (category name)"
		err := SendDeletableMessage(m.Sender, b, text, config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return

	}

	if m.ReplyTo != nil {
		if !lr.RecordExists(int32(m.ReplyTo.ID)) {
			text := "You can not edit this message"
			err := SendDeletableMessage(m.Sender, b, text, config.NotificationTimeout)
			if err != nil {
				logrus.Error(err)
				return
			}
		} else {
			logrus.Info("Start editing")

			go DeleteMessage(m.ReplyTo, b, 0)
			err := lr.DeleteRecordsByMessageID(int32(m.ReplyTo.ID))
			if err != nil {
				logrus.Error(err)
			}

			logrus.Info("Remove all related records")

			SaveItems(items, int32(m.ID), m.Sender, b, lr, config)
		}

	} else {
		SaveItems(items, int32(m.ID), m.Sender, b, lr, config)
	}
}

// HandleEdit allow to edit infromation from db for following message
func HandleEdit(m *tb.Message, b *tb.Bot, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, b, config) {
		return
	}

	go Notify(m.Sender, b, tb.Typing)

	item := GetItem(m.Text)
	logrus.Info("Parsed data", item)

	editLogs(int32(m.ID), m.Sender, b, item, lr, config)
}

// editLogs tries to edit log
func editLogs(messageID int32, sender *tb.User, b *tb.Bot, items []Item, lr LogItemRepository, config Config) {
	logrus.Info("Start editing")
	var text string
	var err error

	if len(items) == 0 {
		text = "Use the following format: `item amount`. *For example*: tea 10 (category name)"
		err = SendDeletableMessage(sender, b, text, config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return

	}

	err = lr.DeleteRecordsByMessageID(messageID)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("Remove all related records")

	SaveItems(items, messageID, sender, b, lr, config)
}

// HandleDelete allow to delete infromation from db for following message
func HandleDelete(m *tb.Message, b *tb.Bot, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleDelete request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, b, config) {
		return
	}

	if m.ReplyTo == nil {
		text := "You should reply for a message which you want to delete ↩️"
		err := SendDeletableMessage(m.Sender, b, text, config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
			return
		}
	} else {
		err := lr.DeleteRecordsByMessageID(int32(m.ReplyTo.ID))
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Info("Remove all related records")

		text := "`Remove item`"
		err = SendDeletableMessage(m.Sender, b, text, config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
			return
		}

		go DeleteMessage(m, b, config.NotificationTimeout)
		if m.ReplyTo != nil {
			go DeleteMessage(m.ReplyTo, b, config.NotificationTimeout)
		}

	}
}

// HandleStatsAllByMonth allow to get information grouped by months
func HandleStatsAllByMonth(m *tb.Message, b *tb.Bot, c stats.StatsClient, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleStatsAllByMonth request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, b, config) {
		return
	}

	go Notify(m.Sender, b, tb.UploadingDocument)

	items, err := lr.GetAggregatedRecords()
	if err != nil {
		logrus.Error(err)
	}

	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		err = SendDeletableMessage(m.Sender, b, "There are not any records yet 😒", config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	itemsForAnalyze := PrepareForAnalyze(items)

	logrus.Info("Call GetMonthAmountStat")
	monthAmountStat, err := c.GetMonthAmountStat(context.Background(), &stats.LogItemQueryMessage{
		LogMessagesAggregated: itemsForAnalyze,
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	fileName := fmt.Sprintf("%v-%v-stat.png", m.Sender.ID, Timestamp())
	err = SendDocumentFromReader(m.Sender, b, fileName, monthAmountStat.Res, config)
	if err != nil {
		logrus.Error(err)
	}

	go Notify(m.Sender, b, tb.UploadingPhoto)

	logrus.Info("Call GetMonthStat")
	monthStat, err := c.GetMonthStat(context.Background(), &stats.LogItemQueryMessage{
		LogMessagesAggregated: itemsForAnalyze,
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("Call GetCategoryStat")
	categoryStat, err := c.GetCategoryStat(context.Background(), &stats.LogItemQueryMessage{
		LogMessagesAggregated: itemsForAnalyze,
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	categoryStatDocument := &tb.Photo{File: tb.FromReader(bytes.NewReader(categoryStat.Res))}
	monthStatDocument := &tb.Photo{File: tb.FromReader(bytes.NewReader(monthStat.Res))}

	err = SendAlbum(m.Sender, b, tb.Album{monthStatDocument, categoryStatDocument}, config)
	if err != nil {
		logrus.Error(err)
	}

	go DeleteMessage(m, b, config.NotificationTimeout)

}

// HandleStatsByCategoryForCurrentMonth allow to get information grouped by categories for current month
func HandleStatsByCategoryForCurrentMonth(m *tb.Message, b *tb.Bot, c stats.StatsClient, lr LogItemRepository, config Config) {
	logrus.Infof("Start HandleStatsByCategoryForCurrentMonth request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, b, config) {
		return
	}

	go Notify(m.Sender, b, tb.UploadingPhoto)

	items, err := lr.GetAggregatedRecordsCurrentMonth()
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		err = SendDeletableMessage(m.Sender, b, "There are not any records yet 😒", config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	logrus.Info("Call GetCategoryStat")
	stat, err := c.GetCategoryStat(context.Background(), &stats.LogItemQueryMessage{
		LogMessagesAggregated: PrepareForAnalyze(items),
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	photo := &tb.Photo{File: tb.FromReader(bytes.NewReader(stat.Res))}

	err = SendDeletableMessage(m.Sender, b, photo, config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
	}

	go DeleteMessage(m, b, config.NotificationTimeout)
}

// HandleExport allow to export data into csv file
func HandleExport(m *tb.Message, b *tb.Bot, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, b, config) {
		return
	}

	go Notify(m.Sender, b, tb.UploadingDocument)

	items, err := lr.GetRecords()
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		err = SendDeletableMessage(m.Sender, b, "There are not any records yet 😒", config.NotificationTimeout)
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

	err = SendDeletableMessage(m.Sender, b, document, config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Send file to ", m.Sender.ID)

	go DeleteMessage(m, b, config.NotificationTimeout)
}

// HandleIntegration allows to add new integration for example bank
func HandleIntegration(items <-chan Item, b *tb.Bot, lr LogItemRepository, config Config) {
	recipient := User{
		ID: config.ChatID,
	}
	var err error
	for item := range items {
		if item.IsValid() {
			if item.Category == "" {
				item.Category, err = lr.FetchMostRelevantCategory(item.Name)
				if err != nil {
					logrus.Error(err)
				}
			}

			text := fmt.Sprintf("%s %.2f %s", item.Name, item.Amount, item.Category)
			message, err := SendMessage(recipient, b, text)

			if err != nil {
				logrus.Error(err)
				continue
			}

			_, err = item.ProcessSaving(int32(message.ID), recipient.ID, b, lr, config)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}
