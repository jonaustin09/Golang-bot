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

	err = SendServiceMessage(m.Sender, b, text, config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
		return
	}
}

// HandleNewMessage process new messages
func HandleNewMessage(m *tb.Message, b *tb.Bot, inputLogRepository InputLogRepository, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleNewMessage request with %s by %v", m.Text, m.Sender.ID)
	parsedData := GetParsedData(m.Text)
	logrus.Info("Parsed data", parsedData)

	err := inputLogRepository.CreateRecord(m.Text, int32(m.Sender.ID))
	if err != nil {
		logrus.Error(err)
	}

	if m.ReplyTo != nil {
		if !lr.RecordExists(int32(m.ReplyTo.ID)) {
			text := "You can not edit this message"
			err := SendServiceMessage(m.Sender, b, text, config.NotificationTimeout)
			if err != nil {
				logrus.Error(err)
				return
			}
		} else {
			editLogs(int32(m.ReplyTo.ID), m.Sender, b, parsedData, lr, config)
		}

	} else {
		var text string

		if !(len(parsedData) > 0) {
			text = "Use the following format: `item amount`. *For example*: tea 10 (repository name)"
			err := SendServiceMessage(m.Sender, b, text, config.NotificationTimeout)
			if err != nil {
				logrus.Error(err)
				return
			}

		} else {
			for _, item := range parsedData {
				if item.Category == "" {
					item.Category, err = lr.FetchMostRelevantCategory(item.Name, int32(m.Sender.ID))
					if err != nil {
						logrus.Error(err)
					}
				}
				logItem, err := lr.CreateRecord(item, int32(m.ID), int32(m.Sender.ID))
				if err != nil {
					logrus.Error(err)
					return
				}

				text = fmt.Sprintf("`Saved: %s`", logItem)
				err = SendServiceMessage(m.Sender, b, text, config.NotificationTimeout)
				if err != nil {
					logrus.Error(err)
				}
			}

		}
	}
}

// HandleEdit allow to edit infromation from db for following message
func HandleEdit(m *tb.Message, b *tb.Bot, inputLogRepository InputLogRepository, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)

	parsedData := GetParsedData(m.Text)
	logrus.Info("Parsed data", parsedData)

	err := inputLogRepository.CreateRecord(m.Text, int32(m.Sender.ID))
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

	if !(len(parsedData) > 0) {
		text = "Use the following format: `item amount`. *For example*: tea 10 (repository name)"
		err = SendServiceMessage(sender, b, text, config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}

	} else {
		text := "`Remove related items`"
		err = SendServiceMessage(sender, b, text, config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
		}

		err := lr.DeleteRecordsByMessageID(messageID)
		if err != nil {
			logrus.Error(err)
		}

		logrus.Info("Remove all related records")

		for _, item := range parsedData {
			logItem, err := lr.CreateRecord(item, int32(messageID), int32(sender.ID))
			if err != nil {
				logrus.Error(err)
			}

			text = fmt.Sprintf("`Create: %s`", logItem.String())
			err = SendServiceMessage(sender, b, text, config.NotificationTimeout)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}

// HandleDelete allow to delete infromation from db for following message
func HandleDelete(m *tb.Message, b *tb.Bot, lr LogItemRepository, config Config) {
	logrus.Infof("Start handleDelete request with %s by %v", m.Text, m.Sender.ID)
	if m.ReplyTo == nil {
		text := "You should reply for a message which you want to delete ↩️"
		err := SendServiceMessage(m.Sender, b, text, config.NotificationTimeout)
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
		err = SendServiceMessage(m.Sender, b, text, config.NotificationTimeout)
		if err != nil {
			logrus.Error(err)
			return
		}
	}
}

// HandleStatsAllByMonth allow to get information grouped by monthes
func HandleStatsAllByMonth(m *tb.Message, b *tb.Bot, c stats.StatsClient, lr LogItemRepository) {
	logrus.Infof("Start handleStatsAllByMonth request with %s by %v", m.Text, m.Sender.ID)
	var err error

	items, err := lr.GetRecordsByTelegramID(int32(m.Sender.ID))
	if err != nil {
		logrus.Error(err)
	}

	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		_, err = b.Send(m.Sender, "There are not any records yet 😒")
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

	monthAmountStatDocument := &tb.Photo{File: tb.FromReader(bytes.NewReader(monthAmountStat.Res))}
	monthStatDocument := &tb.Photo{File: tb.FromReader(bytes.NewReader(monthStat.Res))}

	_, err = b.SendAlbum(m.Sender, tb.Album{monthStatDocument, monthAmountStatDocument})
	if err != nil {
		logrus.Error(err)
	}

}

// HandleStatsByCategory allow to get information grouped by categories
func HandleStatsByCategory(m *tb.Message, b *tb.Bot, c stats.StatsClient, lr LogItemRepository) {
	logrus.Infof("Start handleStatsAllByMonth request with %s by %v", m.Text, m.Sender.ID)
	var err error

	items, err := lr.GetRecordsByTelegramID(int32(m.Sender.ID))
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		_, err = b.Send(m.Sender, "There are not any records yet 😒")
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	album := tb.Album{}

	logrus.Info("Call GetCategoryStat")
	statAll, err := c.GetCategoryStat(context.Background(), &stats.LogItemQueryMessage{
		LogItems: PrepareForAnalyze(items),
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	album = append(album, &tb.Photo{File: tb.FromReader(bytes.NewReader(statAll.Res))})

	items, err = lr.GetRecordsByTelegramIDCurrentMonth(int32(m.Sender.ID))
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("Fetch items count %v", len(items))

	if len(items) != 0 {
		logrus.Info("Call GetCategoryStat")
		statByCurrentMonth, err := c.GetCategoryStat(context.Background(), &stats.LogItemQueryMessage{
			LogItems: PrepareForAnalyze(items),
		})
		if err != nil {
			logrus.Error(err)
			return
		}

		album = append(tb.Album{&tb.Photo{File: tb.FromReader(bytes.NewReader(statByCurrentMonth.Res))}}, album...)
	}

	_, err = b.SendAlbum(m.Sender, album)
	if err != nil {
		logrus.Error(err)
	}
}

// HandleExport allow to export data into csv file
func HandleExport(m *tb.Message, b *tb.Bot, lr LogItemRepository) {
	logrus.Infof("Start handleEdit request with %s by %v", m.Text, m.Sender.ID)
	var err error

	items, err := lr.GetRecordsByTelegramID(int32(m.Sender.ID))
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("Fetch items count %v", len(items))

	if len(items) == 0 {
		_, err = b.Send(m.Sender, "There are not any records yet 😒")
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
			return
		}
	}()
	defer func() {
		err := file.Close()
		if err != nil {
			logrus.Error(err)
			return
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

	_, err = b.Send(m.Sender, document)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Send file to ", m.Sender.ID)
}
