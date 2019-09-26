package moneybot

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Timestamp returns unix now time
func Timestamp() int32 {
	res := time.Now().UnixNano() / int64(time.Second)
	return int32(res)
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
	message, err := b.Send(to, d, tb.ModeMarkdown, tb.Silent)
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

// SendDocumentFromReader sends bytes as file
func SendDocumentFromReader(to tb.Recipient, b *tb.Bot, fileName string, file []byte, config Config) error {
	err := ioutil.WriteFile(fileName, file, 0644)
	if err != nil {
		return err
	}
	logrus.Info("Create file")

	defer func() {
		err := os.Remove(fileName)
		if err != nil {
			logrus.Error(err)
		}
	}()

	document := &tb.Document{File: tb.FromDisk(fileName)}
	err = SendDeletableMessage(to, b, document, config.NotificationTimeout)

	return err
}

// SendAlbum sends album and remove message after timeout
func SendAlbum(to tb.Recipient, b *tb.Bot, a tb.Album, config Config) error {
	messages, err := b.SendAlbum(to, a)
	if err != nil {
		return err
	}

	for _, m := range messages {
		tmp := m
		go DeleteMessage(&tmp, b, config.NotificationTimeout)
	}

	return nil
}

// Notify notifies user that action is started
func Notify(to tb.Recipient, b *tb.Bot, action tb.ChatAction) {
	err := b.Notify(to, action)
	if err != nil {
		logrus.Error(err)
	}
}
