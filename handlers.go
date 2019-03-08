package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

func handleStart(m *tb.Message, b *tb.Bot) {
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
	check(err)

}

func handleNewMessage(m *tb.Message, b *tb.Bot) {
	logrus.Infof("Start handleNewMessage request with %s by %v", m.Text, m.Sender.ID)
	parsedData := getParsedData(m.Text)
	logrus.Info("Parsed data", parsedData)

	if m.ReplyTo != nil {
		if !recordExists(uint64(m.ReplyTo.ID)) {
			text := "You can not edit this message"
			err := sendServiceMessage(m.Sender, b, text)
			check(err)
		} else {
			editLogs(uint64(m.ReplyTo.ID), m.Sender, b, parsedData)
		}

	} else {
		var text string

		if !parsedDataIsValid(parsedData) {
			text = "Use the following format: `item amount`. *For example*: tea 10 (category name)"
			err := sendServiceMessage(m.Sender, b, text)
			check(err)

		} else {
			for _, item := range parsedData {
				logItem := LogItem{}
				err := logItem.createRecord(item, uint64(m.ID), uint64(m.Sender.ID))
				check(err)

				text = fmt.Sprintf("`Saved: %s`", logItem.String())
				err = sendServiceMessage(m.Sender, b, text)
				check(err)
			}

		}
	}

}

func handleEdit(m *tb.Message, b *tb.Bot) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)

	parsedData := getParsedData(m.Text)
	logrus.Info("Parsed data", parsedData)

	editLogs(uint64(m.ID), m.Sender, b, parsedData)
}

func handleDelete(m *tb.Message, b *tb.Bot) {
	logrus.Infof("Start handleDelete request with %s by %v", m.Text, m.Sender.ID)
	if m.ReplyTo == nil {
		text := "You should reply for a message which you want to delete ↩️"
		err := sendServiceMessage(m.Sender, b, text)
		check(err)
	} else {
		err := deleteRecordsByMessageID(uint64(m.ReplyTo.ID))
		check(err)

		text := "Remove item"
		err = sendServiceMessage(m.Sender, b, text)
		check(err)
	}
}

func handleExport(m *tb.Message, b *tb.Bot) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)
	var err error
	items, err := getRecordsByTelegramID(uint64(m.Sender.ID))
	check(err)
	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		_, err = b.Send(m.Sender, "There are not any records yet 😒")
		check(err)
		return
	}

	fileName := fmt.Sprintf("%v-%v-export.csv", m.Sender.ID, timestamp())
	file, err := os.Create(fileName)
	check(err)
	logrus.Info("Create file")

	defer os.Remove(fileName)
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, item := range items {
		err = writer.Write(item.toCSV())
		check(err)
	}
	writer.Flush()
	logrus.Info("Save file")

	document := &tb.Document{File: tb.FromDisk(fileName)}

	_, err = b.Send(m.Sender, document)
	check(err)
	logrus.Info("Send file to ", m.Sender.ID)
}

func editLogs(messageID uint64, sender *tb.User, b *tb.Bot, parsedData []ParsedData) {
	logrus.Info("Start editing")
	var text string
	var err error

	if !parsedDataIsValid(parsedData) {
		text = "Use the following format: `item amount`. *For example*: tea 10 (category name)"
		err = sendServiceMessage(sender, b, text)
		check(err)

	} else {
		text := "`Remove related items`"
		err = sendServiceMessage(sender, b, text)
		check(err)

		err := deleteRecordsByMessageID(messageID)
		check(err)
		logrus.Info("Remove all related records")

		for _, item := range parsedData {
			logItem := LogItem{}
			err = logItem.createRecord(item, uint64(messageID), uint64(sender.ID))
			check(err)

			text = fmt.Sprintf("`Create: %s`", logItem.String())
			err = sendServiceMessage(sender, b, text)
			check(err)
		}
	}
}
