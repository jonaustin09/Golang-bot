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
func HandleStart(m *tb.Message, b *tb.Bot, ur UserRepository, config Config) {
	logrus.Infof("Start handleStart request with %s by %v", m.Text, m.Sender.ID)
	var text string

	uid := int32(m.Sender.ID)

	_, err := ur.FetchOrCreate(uid, m.Sender.FirstName, m.Sender.LastName, m.Sender.LanguageCode, m.Sender.Username)
	if err != nil {
		logrus.Error(err)
		return
	}
	text = "Hello there i'll help you with your finances! \n" +
		"Use the following format: `item amount`. *For example*: tea 10 (repository name)"

	err = SendMessage(m.Sender, b, text, config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
		return
	}
}

// HandleNewMessage process new messages
func HandleNewMessage(m *tb.Message, b *tb.Bot, inputLogRepository InputLogRepository, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleNewMessage request with %s by %v", m.Text, m.Sender.ID)

	err := b.Notify(m.Sender, tb.Typing)
	if err != nil {
		logrus.Error(err)
	}

	parsedData := GetParsedData(m.Text)
	logrus.Info("Parsed data", parsedData)

	err = inputLogRepository.CreateRecord(m.Text, int32(m.Sender.ID))
	if err != nil {
		logrus.Error(err)
	}

	if m.ReplyTo != nil {
		if !lr.RecordExists(int32(m.ReplyTo.ID)) {
			text := "You can not edit this message"
			err := SendMessage(m.Sender, b, text, config.NotificationTimeout)
			if err != nil {
				logrus.Error(err)
				return
			}
		} else {
			editLogs(int32(m.ReplyTo.ID), m.Sender, b, parsedData, lr, config)
		}

	} else {
		var text string

		if len(parsedData) == 0 {
			text = "Use the following format: `item amount`. *For example*: tea 10 (repository name)"
			err := SendMessage(m.Sender, b, text, config.NotificationTimeout)
			if err != nil {
				logrus.Error(err)
				return
			}
			return

		}

		SaveParsedData(parsedData, int32(m.ID), m.Sender, b, lr, config)
	}
}

// HandleEdit allow to edit infromation from db for following message
func HandleEdit(m *tb.Message, b *tb.Bot, inputLogRepository InputLogRepository, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)

	err := b.Notify(m.Sender, tb.Typing)
	if err != nil {
		logrus.Error(err)
	}

	parsedData := GetParsedData(m.Text)
	logrus.Info("Parsed data", parsedData)

	err = inputLogRepository.CreateRecord(m.Text, int32(m.Sender.ID))
	if err != nil {
		logrus.Error(err)
	}

	editLogs(int32(m.ID), m.Sender, b, parsedData, lr, config)
}

// editLogs tries to edit log
func editLogs(messageID int32, sender *tb.User, b *tb.Bot, parsedData []ParsedData, lr LogItemRepository, config Config) {
	logrus.Info("Start editing")
	var text string
	var err error

	if len(parsedData) == 0 {
		text = "Use the following format: `item amount`. *For example*: tea 10 (repository name)"
		err = SendMessage(sender, b, text, config.NotificationTimeout)
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

	SaveParsedData(parsedData, messageID, sender, b, lr, config)
}

// HandleDelete allow to delete infromation from db for following message
func HandleDelete(m *tb.Message, b *tb.Bot, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleDelete request with %s by %v", m.Text, m.Sender.ID)
	if m.ReplyTo == nil {
		text := "You should reply for a message which you want to delete ↩️"
		err := SendMessage(m.Sender, b, text, config.NotificationTimeout)
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
		err = SendMessage(m.Sender, b, text, config.NotificationTimeout)
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

// HandleStatsAllByMonth allow to get information grouped by monthes
func HandleStatsAllByMonth(m *tb.Message, b *tb.Bot, c stats.StatsClient, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleStatsAllByMonth request with %s by %v", m.Text, m.Sender.ID)

	err := b.Notify(m.Sender, tb.UploadingPhoto)
	if err != nil {
		logrus.Error(err)
	}

	items, err := lr.GetRecordsByTelegramID(int32(m.Sender.ID))
	if err != nil {
		logrus.Error(err)
	}

	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		err = SendMessage(m.Sender, b, "There are not any records yet 😒", config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	itemsForAnalyze := PrepareForAnalyze(items)

	logrus.Info("Call GetMonthStat")
	monthStat, err := c.GetMonthStat(context.Background(), &stats.LogItemQueryMessage{
		LogItems: itemsForAnalyze,
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("Call GetMonthAmountStat")
	monthAmountStat, err := c.GetMonthAmountStat(context.Background(), &stats.LogItemQueryMessage{
		LogItems: itemsForAnalyze,
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	var fileName string

	fileName = fmt.Sprintf("%v-%v-stat.png", m.Sender.ID, Timestamp())
	err = SendDocumentFromReader(m.Sender, b, fileName, monthAmountStat.Res, config)
	if err != nil {
		logrus.Error(err)
	}

	fileName = fmt.Sprintf("%v-%v-stat.png", m.Sender.ID, Timestamp())
	err = SendDocumentFromReader(m.Sender, b, fileName, monthStat.Res, config)
	if err != nil {
		logrus.Error(err)
	}

	go DeleteMessage(m, b, config.NotificationTimeout)

}

// HandleStatsByCategory allow to get information grouped by categories
func HandleStatsByCategory(m *tb.Message, b *tb.Bot, c stats.StatsClient, lr LogItemRepository, config Config) {
	logrus.Infof("Start HandleStatsByCategory request with %s by %v", m.Text, m.Sender.ID)
	err := b.Notify(m.Sender, tb.UploadingPhoto)
	if err != nil {
		logrus.Error(err)
	}

	items, err := lr.GetRecordsByTelegramID(int32(m.Sender.ID))
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		err = SendMessage(m.Sender, b, "There are not any records yet 😒", config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	logrus.Info("Call GetCategoryStat")
	stat, err := c.GetCategoryStat(context.Background(), &stats.LogItemQueryMessage{
		LogItems: PrepareForAnalyze(items),
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	photo := &tb.Photo{File: tb.FromReader(bytes.NewReader(stat.Res))}

	err = SendMessage(m.Sender, b, photo, config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
	}

	go DeleteMessage(m, b, config.NotificationTimeout)
}

// HandleStatsByCategoryForCurrentMonth allow to get information grouped by categories for current month
func HandleStatsByCategoryForCurrentMonth(m *tb.Message, b *tb.Bot, c stats.StatsClient, lr LogItemRepository, config Config) {
	logrus.Infof("Start HandleStatsByCategoryForCurrentMonth request with %s by %v", m.Text, m.Sender.ID)
	err := b.Notify(m.Sender, tb.UploadingPhoto)
	if err != nil {
		logrus.Error(err)
	}

	items, err := lr.GetRecordsByTelegramIDCurrentMonth(int32(m.Sender.ID))
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		err = SendMessage(m.Sender, b, "There are not any records yet 😒", config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	logrus.Info("Call GetCategoryStat")
	stat, err := c.GetCategoryStat(context.Background(), &stats.LogItemQueryMessage{
		LogItems: PrepareForAnalyze(items),
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	photo := &tb.Photo{File: tb.FromReader(bytes.NewReader(stat.Res))}

	err = SendMessage(m.Sender, b, photo, config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
	}

	go DeleteMessage(m, b, config.NotificationTimeout)
}

// HandleExport allow to export data into csv file
func HandleExport(m *tb.Message, b *tb.Bot, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)
	err := b.Notify(m.Sender, tb.UploadingDocument)
	if err != nil {
		logrus.Error(err)
	}

	items, err := lr.GetRecordsByTelegramID(int32(m.Sender.ID))
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		err = SendMessage(m.Sender, b, "There are not any records yet 😒", config.NotificationTimeout)
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

	err = SendMessage(m.Sender, b, document, config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Send file to ", m.Sender.ID)

	go DeleteMessage(m, b, config.NotificationTimeout)
}
