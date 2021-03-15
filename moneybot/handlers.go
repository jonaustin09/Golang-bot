package moneybot

import (
	"encoding/csv"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
)


func (a Application) handleStart(m *tb.Message) {
	logrus.Infof("Start handleStart request with %s by %v", m.Text, m.Sender.ID)

	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	text := "Hello there i'll help you with your finances! \n" +
		"Use the following format: `item amount`. *For example*: tea 10 (repository name) \n" +
		"To delete message start to replay what you want to delete and type button 'delete'"

	_, err := SendMessage(m.Sender, a.Bot, text)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func (a Application) handleNewMessage(m *tb.Message) {
	logrus.Infof("Start handleNewMessage request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	go Notify(m.Sender, a.Bot, tb.Typing)

	items := GetItem(m.Text)
	logrus.Info("Parsed data", items)

	if len(items) == 0 {
		text := "Use the following format: `item amount`. *For example*: tea 10 (category name)"
		err := SendDeletableMessage(m.Sender, a.Bot, text, a.Config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}

		go DeleteMessage(m, a.Bot, a.Config.NotificationTimeout)

		return

	}

	if m.ReplyTo != nil {
		if !a.LogItemRepository.RecordExists(int32(m.ReplyTo.ID)) {
			text := "You can not edit this message"
			err := SendDeletableMessage(m.Sender, a.Bot, text, a.Config.NotificationTimeout)
			if err != nil {
				logrus.Error(err)
				return
			}
		} else {
			logrus.Info("Start editing")

			go DeleteMessage(m.ReplyTo, a.Bot, 0)
			err := a.LogItemRepository.DeleteRecordsByMessageID(int32(m.ReplyTo.ID))
			if err != nil {
				logrus.Error(err)
			}

			logrus.Info("Remove all related records")

			SaveItems(items, int32(m.ID), m.Sender, a.Bot, a.LogItemRepository, a.Config)
		}

	} else {
		SaveItems(items, int32(m.ID), m.Sender, a.Bot, a.LogItemRepository, a.Config)
	}
}

func (a Application) handleEdit(m *tb.Message) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	go Notify(m.Sender, a.Bot, tb.Typing)

	item := GetItem(m.Text)
	logrus.Info("Parsed data", item)

	a.editLogs(int32(m.ID), m.Sender, item)
}

func (a Application) editLogs(messageID int32, sender *tb.User, items []Item) {
	logrus.Info("Start editing")
	var text string
	var err error

	if len(items) == 0 {
		text = "Use the following format: `item amount`. *For example*: tea 10 (category name)"
		err = SendDeletableMessage(sender, a.Bot, text, a.Config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}
		return

	}

	err = a.LogItemRepository.DeleteRecordsByMessageID(messageID)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("Remove all related records")

	SaveItems(items, messageID, sender, a.Bot, a.LogItemRepository, a.Config)
}

func (a Application) handleDelete(m *tb.Message) {
	logrus.Infof("Start handleDelete request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	if m.ReplyTo == nil {
		text := "You should reply for a message which you want to delete ↩️"
		err := SendDeletableMessage(m.Sender, a.Bot, text, a.Config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
			return
		}
	} else {
		err := a.LogItemRepository.DeleteRecordsByMessageID(int32(m.ReplyTo.ID))
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Info("Remove all related records")

		text := "`Remove item`"
		err = SendDeletableMessage(m.Sender, a.Bot, text, a.Config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
			return
		}

		go DeleteMessage(m, a.Bot, a.Config.NotificationTimeout)
		if m.ReplyTo != nil {
			go DeleteMessage(m.ReplyTo, a.Bot, a.Config.NotificationTimeout)
		}

	}
}


func (a Application) handleExport(m *tb.Message) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)
	if isForbidden(m, a.Bot, a.Config) {
		return
	}

	go Notify(m.Sender, a.Bot, tb.UploadingDocument)

	items, err := a.LogItemRepository.GetRecords()
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		err = SendDeletableMessage(m.Sender, a.Bot, "There are not any records yet 😒", a.Config.NotificationTimeout)
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

	err = SendDeletableMessage(m.Sender, a.Bot, document, a.Config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Send file to ", m.Sender.ID)

	go DeleteMessage(m, a.Bot, a.Config.NotificationTimeout)
}

func (a Application) handleIntegration(items <-chan BankData) {
	var err error
	for event := range items {
		if event.Item.IsValid() {
			if event.Item.Category == "" {
				event.Item.Category, err = a.LogItemRepository.FetchMostRelevantCategory(event.Item.Name)
				if err != nil {
					logrus.Error(err)
				}
			}

			text := fmt.Sprintf("%s %.2f %s", event.Item.Name, event.Item.Amount, event.Item.Category)
			message, err := SendMessage(event.Account, a.Bot, text)

			if err != nil {
				logrus.Error(err)
				continue
			}

			_, err = event.Item.ProcessSaving(int32(message.ID), event.Account.Username, a.LogItemRepository)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}

