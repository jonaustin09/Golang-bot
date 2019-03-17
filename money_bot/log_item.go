package money_bot

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/now"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

// LogItem stores information about each log item
type LogItem struct {
	ID             string
	CreatedAt      uint64
	Name           string
	Amount         float64
	MessageID      uint64
	TelegramUserID uint64
	CategoryID     uint64   `gorm:"DEFAULT:9999"`
	Category       Category `gorm:"auto_preload"`
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

	return fmt.Sprintf("%s %s %.2f %s", timeString, logItem.Name, logItem.Amount, inCategotyString)
}

func (logItem *LogItem) toCSV() []string {
	return []string{
		strconv.FormatInt(int64(logItem.CreatedAt), 10),
		logItem.Name,
		fmt.Sprintf("%.2f", logItem.Amount),
		logItem.Category.Name,
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
			if category.ID == 0 {
				logItem.CategoryID = 9999
			} else {
				logItem.CategoryID = category.ID
			}
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

func getRecordsByTelegramIDCurrentMonth(SenderID uint64) ([]LogItem, error) {
	var items []LogItem

	beginOfMonth := uint64(now.BeginningOfMonth().UnixNano() / int64(time.Second))

	if err := Db.Where("telegram_user_id = ? AND created_at >= ?", SenderID, beginOfMonth).Find(&items).Order("created_at").Error; err != nil {
		return nil, err
	}
	return items, nil
}

type res struct {
	Result float64
}

func getSumByTelegramIDCurrentMonth(SenderID uint64) (float64, error) {
	beginOfMonth := uint64(now.BeginningOfMonth().UnixNano() / int64(time.Second))
	val := res{}
	Db.Where("telegram_user_id = ? AND created_at >= ?", SenderID, beginOfMonth).Select("SUM(amount) as result").Model(&LogItem{}).Scan(&val)
	return val.Result, nil
}

func deleteRecordsByMessageID(MessageID uint64) error {
	return Db.Where("message_id = ?", MessageID).Delete(LogItem{}).Error
}

func recordExists(MessageID uint64) bool {
	return !Db.Where("message_id = ?", MessageID).First(&LogItem{}).RecordNotFound()
}
