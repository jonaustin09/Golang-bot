package main

import (
	"time"

	"github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Timestamp returns unix now time
func Timestamp() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Second))
}

// Check check result for errors
func Check(err error) {
	if err != nil {
		logrus.Panic(err)
	}
}

func deleteSystemMessage(m *tb.Message, b *tb.Bot) {
	time.Sleep(NOTIFICATIONTIMEOUT)
	err := b.Delete(m)
	logrus.Info("Remove service message", m)
	Check(err)
}

func sendServiceMessage(to tb.Recipient, b *tb.Bot, text string) error {
	serviceMessage, err := b.Send(to, text, tb.ModeMarkdown, tb.Silent)
	if err != nil {
		return err
	}
	logrus.Info("Send service message", serviceMessage)
	go deleteSystemMessage(serviceMessage, b)

	return nil
}
