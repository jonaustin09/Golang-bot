package moneybot

import (
	"time"

	"github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Timestamp returns unix now time
func timestamp() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Second))
}

func getLocalTime(timestamp uint64) time.Time {
	unixTime := time.Unix(int64(timestamp), 0)
	t, err := time.LoadLocation("Europe/Kiev")
	Check(err)
	return unixTime.In(t)
}

// Check Check result for errors
func Check(err error) {
	if err != nil {
		logrus.Panic(err)
	}
}

func deleteSystemMessage(m *tb.Message, b *tb.Bot) {
	time.Sleep(Confg.NotificationTimeout)
	err := b.Delete(m)
	if err != nil {
		logrus.Info(err)
	}
	logrus.Info("Remove service message ", m.ID)
}

func sendServiceMessage(to tb.Recipient, b *tb.Bot, text string) error {
	serviceMessage, err := b.Send(to, text, tb.ModeMarkdown, tb.Silent)
	if err != nil {
		return err
	}
	logrus.Infof("Send service message %v with text: %s", serviceMessage.ID, serviceMessage.Text)
	go deleteSystemMessage(serviceMessage, b)

	return nil
}
