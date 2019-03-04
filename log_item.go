package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// LogItem stores information about each log item
type LogItem struct {
	ID             string `gorm:"primary_key"`
	CreatedAt      uint64
	Name           string
	Amount         float64
	MessageID      uint64
	TelegramUserID uint64
	CategoryID     uint64
}

func (logItem *LogItem) createRecord(parsedData ParsedData, messageID uint64, senderID uint64) error {
	uid, err := uuid.NewV4()
	Check(err)

	logItem.ID = uid.String()
	logItem.Name = parsedData.Name
	logItem.Amount = parsedData.Amount
	logItem.MessageID = messageID
	logItem.CreatedAt = Timestamp()
	logItem.TelegramUserID = senderID

	if parsedData.HasCategory() {
		category := Category{}
		err = category.fetchOrCreate(parsedData.Category, senderID)
		if err == nil {
			logItem.CategoryID = category.ID
		}

	}

	logrus.Info("Create record logItem ", logItem)
	err = db.Create(&logItem).Error
	return err
}

func (logItem *LogItem) String() string {
	unixTime := time.Unix(int64(logItem.CreatedAt), 0)
	timeString := unixTime.Format("02-01-2006")

	inCategotyString := ""

	category := Category{}
	err := category.fetchByID(logItem.CategoryID)
	if err == nil {
		inCategotyString = fmt.Sprintf("in %s", category.Name)
	}

	return fmt.Sprintf("%s %s %v %s", timeString, logItem.Name, logItem.Amount, inCategotyString)
}

func (logItem *LogItem) toCSV() []string {
	category := Category{}
	category.fetchByID(logItem.CategoryID) // nolint: gosec
	// TODO: think about this extra query

	return []string{
		strconv.FormatInt(int64(logItem.CreatedAt), 10),
		logItem.Name,
		fmt.Sprintf("%f", logItem.Amount),
		category.Name,
	}
}

func (logItem *LogItem) getByMessageID(messageID uint64) error {
	return db.Where("message_id = ?", messageID).First(logItem).Error
}

func (logItem *LogItem) updateRecord(parsedData ParsedData) error {
	if logItem.ID == "" {
		return errors.New("can update only created items")
	}

	logItem.Name = parsedData.Name
	logItem.Amount = parsedData.Amount

	if err := db.Save(&logItem).Error; err != nil {
		return err
	}
	logrus.Info("Update record logItem ", logItem)

	return nil
}

func getRecordsByTelegramID(SenderID uint64) ([]LogItem, error) {
	var items []LogItem
	if err := db.Where("telegram_user_id = ?", SenderID).Find(&items).Order("created_at").Error; err != nil {
		return nil, err
	}
	return items, nil
}
