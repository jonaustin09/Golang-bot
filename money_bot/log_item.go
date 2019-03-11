package money_bot

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
	CategoryID     uint64 `gorm:"DEFAULT:9999"`
}

func (logItem *LogItem) categoryName() string {
	category := Category{}
	category.fetchByID(logItem.CategoryID) // nolint: gosec
	return category.Name
}

func (logItem *LogItem) String() string {
	unixTime := time.Unix(int64(logItem.CreatedAt), 0)
	timeString := unixTime.Format("02-01-2006")

	inCategotyString := ""

	category := Category{}
	err := category.fetchByID(logItem.CategoryID)
	if err != nil {
		err = category.getDefault()
		Check(err)
	}

	inCategotyString = fmt.Sprintf("in %s", category.Name)

	return fmt.Sprintf("%s %s %v %s", timeString, logItem.Name, logItem.Amount, inCategotyString)
}

func (logItem *LogItem) toCSV() []string {
	// TODO: think about this extra query

	return []string{
		strconv.FormatInt(int64(logItem.CreatedAt), 10),
		logItem.Name,
		fmt.Sprintf("%f", logItem.Amount),
		logItem.categoryName(),
	}
}

func (logItem *LogItem) createRecord(parsedData ParsedData, MessageID uint64, senderID uint64) error {
	uid, err := uuid.NewV4()
	Check(err)

	logItem.ID = uid.String()
	logItem.Name = parsedData.Name
	logItem.Amount = parsedData.Amount
	logItem.MessageID = MessageID
	logItem.CreatedAt = timestamp()
	logItem.TelegramUserID = senderID

	if parsedData.hasCategory() {
		category := Category{}
		err = category.fetchOrCreate(parsedData.Category, senderID)
		if err == nil {
			logItem.CategoryID = category.ID
		}
	} else {
		category := Category{}
		err = category.fetchMostRelevantForItem(logItem.Name, logItem.TelegramUserID)
		if err == nil {
			logItem.CategoryID = category.ID
		}
	}

	logrus.Info("Create record logItem ", logItem)
	err = Db.Create(&logItem).Error
	return err
}

func (logItem *LogItem) updateRecord(parsedData ParsedData, senderID uint64) error {
	if logItem.ID == "" {
		return errors.New("can update only created items")
	}

	logItem.Name = parsedData.Name
	logItem.Amount = parsedData.Amount

	if parsedData.hasCategory() {
		category := Category{}
		err := category.fetchOrCreate(parsedData.Category, senderID)
		if err == nil {
			logItem.CategoryID = category.ID
		}
	} else {
		if logItem.CategoryID == 0 {
			category := Category{}
			err := category.fetchMostRelevantForItem(logItem.Name, logItem.TelegramUserID)
			if err == nil {
				logItem.CategoryID = category.ID
			}
		}
	}

	if err := Db.Save(&logItem).Error; err != nil {
		return err
	}
	logrus.Info("Update record logItem ", logItem)

	return nil
}

func getRecordsByTelegramID(SenderID uint64) ([]LogItem, error) {
	var items []LogItem
	if err := Db.Where("telegram_user_id = ?", SenderID).Find(&items).Order("created_at").Error; err != nil {
		return nil, err
	}
	return items, nil
}

func deleteRecordsByMessageID(MessageID uint64) error {
	return Db.Where("message_id = ?", MessageID).Delete(LogItem{}).Error
}

func recordExists(MessageID uint64) bool {
	return !Db.Where("message_id = ?", MessageID).First(&LogItem{}).RecordNotFound()
}
