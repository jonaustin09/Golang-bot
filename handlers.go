package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

func handleStart(m *tb.Message, b *tb.Bot) {
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

	_, err := b.Send(m.Sender, text)
	if err != nil {
		logrus.Panic(err)
	}

}

func handleNewMessage(m *tb.Message, b *tb.Bot) {
	parsedData := GetParsedData(m.Text)
	var text string
	if !parsedData.IsValid() {
		text = "Use the following format: `item amount`. *For example*: tea 10 (category name)"

	} else {
		logItem := LogItem{}
		err := logItem.createRecord(parsedData, uint64(m.ID), uint64(m.Sender.ID))
		if err != nil {
			logrus.Panic(err)
		}

		text = fmt.Sprintf("`Saved: %s`", logItem.String())

	}
	err := sendServiceMessage(m.Sender, b, text)
	Check(err)
}

func handleEdit(m *tb.Message, b *tb.Bot) {
	parsedData := GetParsedData(m.Text)
	var text string
	var err error
	if !parsedData.IsValid() {
		text = "Use the following format: `item amount`. *For example*: tea 10 (category name)"

	} else {
		logItem := LogItem{}

		err = logItem.getByMessageID(uint64(m.ID))
		Check(err)

		err = logItem.updateRecord(parsedData)
		Check(err)

		text = fmt.Sprintf("`Updated: %s`", logItem.String())

	}

	serviceMessage, err := b.Send(m.Sender, text, tb.ModeMarkdown)
	if err != nil {
		logrus.Panic(err)
	}

	go deleteSystemMessage(serviceMessage, b)

}

func handleExport(m *tb.Message, b *tb.Bot) {
	var err error
	items, err := getRecordsByTelegramID(uint64(m.Sender.ID))
	Check(err)

	if len(items) == 0 {
		_, err = b.Send(m.Sender, "There are not any records yet 😒")
		Check(err)
		return
	}

	fileName := fmt.Sprintf("%v-%v-export.csv", m.Sender.ID, Timestamp())
	file, err := os.Create(fileName)
	Check(err)

	defer os.Remove(fileName)
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, item := range items {
		err = writer.Write(item.toCSV())
		Check(err)
	}
	writer.Flush()

	document := &tb.Document{File: tb.FromDisk(fileName)}

	_, err = b.Send(m.Sender, document)
	Check(err)
}
