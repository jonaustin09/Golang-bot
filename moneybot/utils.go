package moneybot

import (
	"time"

	"github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Timestamp returns unix now time
func Timestamp() int32 {
	return int32(time.Now().Unix())
}

// GetLocalTime utc time to Europe/Kiev
func GetLocalTime(timestamp int32) time.Time {
	// TODO: refactor this create user settings
	unixTime := time.Unix(int64(timestamp), 0)
	t, err := time.LoadLocation("Europe/Kiev")
	if err != nil {
		logrus.Panic(err)
	}
	return unixTime.In(t)
}

// DeleteMessage tries to delete sent message
func DeleteMessage(m *tb.Message, b *tb.Bot, timeout time.Duration) {
	time.Sleep(timeout)
	err := b.Delete(m)
	if err != nil {
		logrus.Info(err)
	}
	logrus.Info("Remove message ", m.ID)
}

// SendMessage tries to sent message
func SendMessage(to tb.Recipient, b *tb.Bot, d interface{}) (*tb.Message, error) {
	message, err := b.Send(to, d, tb.Silent)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Send message %v", message.ID, message.Text)
	return message, nil
}

// SendDeletableMessage tries to sent message, will delete after timeout
func SendDeletableMessage(to tb.Recipient, b *tb.Bot, d interface{}, displayTimeout time.Duration) error {
	serviceMessage, err := SendMessage(to, b, d)
	if err != nil {
		return err
	}
	go DeleteMessage(serviceMessage, b, displayTimeout)

	return nil
}

// Notify notifies user that action is started
func Notify(to tb.Recipient, b *tb.Bot, action tb.ChatAction) {
	err := b.Notify(to, action)
	if err != nil {
		logrus.Error(err)
	}
}

func isForbidden(m *tb.Message, b *tb.Bot, config Config) bool {
	if m.Sender.ID != config.SenderID1 && m.Sender.ID != config.SenderID2 {
		text := "You can't use this bot. You should deploy it."
		err := SendDeletableMessage(m.Sender, b, text, config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return true
	}
	return false
}
