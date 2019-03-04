package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

func handleStart(m *tb.Message, b *tb.Bot) {
	logrus.Info("Start handleStart request with", m)
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

func handleNewMessage(m *tb.Message, b *tb.Bot) {
	logrus.Info("Start handleNewMessage request with", m)
	parsedData := GetParsedData(m.Text)
	logrus.Info("Parsed data", parsedData)
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

func handleEdit(m *tb.Message, b *tb.Bot) {
	logrus.Info("Start handleEdit request with", m)
	parsedData := GetParsedData(m.Text)
	logrus.Info("Parsed data", parsedData)

	var text string
	var err error

	if !parsedDataIsValid(parsedData) {
		text = "Use the following format: `item amount`. *For example*: tea 10 (category name)"
		err = sendServiceMessage(m.Sender, b, text)
		Check(err)

	} else {
		for _, item := range parsedData {
			logItem := LogItem{}

			err = logItem.getByMessageID(uint64(m.ID))
			Check(err)

			err = logItem.updateRecord(item)
			Check(err)

			text = fmt.Sprintf("`Updated: %s`", logItem.String())
			err = sendServiceMessage(m.Sender, b, text)
			Check(err)
		}
	}

}

func handleExport(m *tb.Message, b *tb.Bot) {
	logrus.Info("Start handleExport request with", m)
	var err error
	items, err := getRecordsByTelegramID(uint64(m.Sender.ID))
	Check(err)
	logrus.Info("Fetch items", items)

	if len(items) == 0 {
		_, err = b.Send(m.Sender, "There are not any records yet 😒")
		Check(err)
		return
	}

	fileName := fmt.Sprintf("%v-%v-export.csv", m.Sender.ID, Timestamp())
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

	message, err := b.Send(m.Sender, document)
	Check(err)
	logrus.Info("Send file", message)
}
