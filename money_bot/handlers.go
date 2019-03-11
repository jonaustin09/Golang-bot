package money_bot

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/dobrovolsky/money_bot/stats"
	"github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

func HandleStart(m *tb.Message, b *tb.Bot) {
	logrus.Infof("Start handleStart request with %s by %v", m.Text, m.Sender.ID)
	var text string

	uid := uint64(m.Sender.ID)

	user := User{}
	isCreated := user.fetchOrCreate(uid, m.Sender.FirstName, m.Sender.LastName,
		m.Sender.LanguageCode, m.Sender.Username)

	if isCreated {
		text = "Hello there i'll help you with your finances!"
	} else {
		text = "Welcome back!"
	}

	err := sendServiceMessage(m.Sender, b, text)
	Check(err)

}

func HandleNewMessage(m *tb.Message, b *tb.Bot) {
	logrus.Infof("Start handleNewMessage request with %s by %v", m.Text, m.Sender.ID)
	parsedData := getParsedData(m.Text)
	logrus.Info("Parsed data", parsedData)

	if m.ReplyTo != nil {
		if !recordExists(uint64(m.ReplyTo.ID)) {
			text := "You can not edit this message"
			err := sendServiceMessage(m.Sender, b, text)
			Check(err)
		} else {
			editLogs(uint64(m.ReplyTo.ID), m.Sender, b, parsedData)
		}

	} else {
		var text string

		if !parsedDataIsValid(parsedData) {
			text = "Use the following format: `item amount`. *For example*: tea 10 (category name)"
			err := sendServiceMessage(m.Sender, b, text)
			Check(err)

		} else {
			for _, item := range parsedData {
				logItem := LogItem{}
				err := logItem.createRecord(item, uint64(m.ID), uint64(m.Sender.ID))
				Check(err)

				text = fmt.Sprintf("`Saved: %s`", logItem.String())
				err = sendServiceMessage(m.Sender, b, text)
				Check(err)
			}

		}
	}

}

func HandleEdit(m *tb.Message, b *tb.Bot) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)

	parsedData := getParsedData(m.Text)
	logrus.Info("Parsed data", parsedData)

	editLogs(uint64(m.ID), m.Sender, b, parsedData)
}

func HandleDelete(m *tb.Message, b *tb.Bot) {
	logrus.Infof("Start handleDelete request with %s by %v", m.Text, m.Sender.ID)
	if m.ReplyTo == nil {
		text := "You should reply for a message which you want to delete ↩️"
		err := sendServiceMessage(m.Sender, b, text)
		Check(err)
	} else {
		err := deleteRecordsByMessageID(uint64(m.ReplyTo.ID))
		Check(err)
		logrus.Info("Remove all related records")

		text := "`Remove item`"
		err = sendServiceMessage(m.Sender, b, text)
		Check(err)
	}
}

func HandleStatsAllByMonth(m *tb.Message, b *tb.Bot, c stats.StatsClient) {
	logrus.Infof("Start handleStatsAllByMonth request with %s by %v", m.Text, m.Sender.ID)
	var err error
	items, err := getRecordsByTelegramID(uint64(m.Sender.ID))
	Check(err)
	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		_, err = b.Send(m.Sender, "There are not any records yet 😒")
		Check(err)
		return
	}

	itemsForAnalyze := make([]*stats.LogItemMessage, 0, len(items))
	for _, item := range items {
		itemsForAnalyze = append(itemsForAnalyze, &stats.LogItemMessage{
			CreatedAt: int64(item.CreatedAt),
			Name:      item.Name,
			Amount:    float32(item.Amount),
			Category:  item.categoryName(),
		})
	}

	response, err := c.GetAllTimeByMonthStat(context.Background(), &stats.LogItemQueryMessage{
		LogItems: itemsForAnalyze,
	})
	Check(err)

	fileName := fmt.Sprintf("%v-%v-stats.png", m.Sender.ID, timestamp())
	file, err := os.Create(fileName)
	Check(err)

	defer os.Remove(fileName)
	defer file.Close()

	file.Write(response.Res)

	document := &tb.Photo{File: tb.FromDisk(fileName)}

	_, err = b.Send(m.Sender, document)
	Check(err)
}

func HandleStatsAllByCategory(m *tb.Message, b *tb.Bot, c stats.StatsClient) {
	logrus.Infof("Start handleStatsAllByMonth request with %s by %v", m.Text, m.Sender.ID)
	var err error
	items, err := getRecordsByTelegramID(uint64(m.Sender.ID))
	Check(err)
	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		_, err = b.Send(m.Sender, "There are not any records yet 😒")
		Check(err)
		return
	}

	itemsForAnalyze := make([]*stats.LogItemMessage, 0, len(items))
	for _, item := range items {
		itemsForAnalyze = append(itemsForAnalyze, &stats.LogItemMessage{
			CreatedAt: int64(item.CreatedAt),
			Name:      item.Name,
			Amount:    float32(item.Amount),
			Category:  item.categoryName(),
		})
	}

	response, err := c.GetAllTimeCategoryStat(context.Background(), &stats.LogItemQueryMessage{
		LogItems: itemsForAnalyze,
	})
	Check(err)

	fileName := fmt.Sprintf("%v-%v-stats.png", m.Sender.ID, timestamp())
	file, err := os.Create(fileName)
	Check(err)

	defer os.Remove(fileName)
	defer file.Close()

	file.Write(response.Res)

	document := &tb.Photo{File: tb.FromDisk(fileName)}

	_, err = b.Send(m.Sender, document)
	Check(err)
}

func HandleExport(m *tb.Message, b *tb.Bot) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)
	var err error
	items, err := getRecordsByTelegramID(uint64(m.Sender.ID))
	Check(err)
	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		_, err = b.Send(m.Sender, "There are not any records yet 😒")
		Check(err)
		return
	}

	fileName := fmt.Sprintf("%v-%v-export.csv", m.Sender.ID, timestamp())
	file, err := os.Create(fileName)
	Check(err)
	logrus.Info("Create file")

	defer os.Remove(fileName)
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, item := range items {
		err = writer.Write(item.toCSV())
		Check(err)
	}
	writer.Flush()
	logrus.Info("Save file")

	document := &tb.Document{File: tb.FromDisk(fileName)}

	_, err = b.Send(m.Sender, document)
	Check(err)
	logrus.Info("Send file to ", m.Sender.ID)
}

func editLogs(messageID uint64, sender *tb.User, b *tb.Bot, parsedData []ParsedData) {
	logrus.Info("Start editing")
	var text string
	var err error

	if !parsedDataIsValid(parsedData) {
		text = "Use the following format: `item amount`. *For example*: tea 10 (category name)"
		err = sendServiceMessage(sender, b, text)
		Check(err)

	} else {
		text := "`Remove related items`"
		err = sendServiceMessage(sender, b, text)
		Check(err)

		err := deleteRecordsByMessageID(messageID)
		Check(err)
		logrus.Info("Remove all related records")

		for _, item := range parsedData {
			logItem := LogItem{}
			err = logItem.createRecord(item, uint64(messageID), uint64(sender.ID))
			Check(err)

			text = fmt.Sprintf("`Create: %s`", logItem.String())
			err = sendServiceMessage(sender, b, text)
			Check(err)
		}
	}
}