package moneybot

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/sirupsen/logrus"
)

// LogItemRepository represent the logItem repository contract
type LogItemRepository interface {
	CreateRecord(parsedData Item, MessageID int32, senderID int32) (*LogItem, error)
	UpdateRecord(logItem *LogItem, parsedData Item, senderID int32) error
	GetRecordsByTelegramID(SenderID int32) ([]LogItem, error)
	GetRecordsByTelegramIDCurrentMonth(SenderID int32) ([]LogItem, error)
	DeleteRecordsByMessageID(MessageID int32) error
	RecordExists(MessageID int32) bool
	FetchMostRelevantCategory(name string, telegramUserID int32) (string, error)
}

// NewGormLogItemRepository creates new repository
func NewGormLogItemRepository(db *gorm.DB) LogItemRepository {
	return GormLogItemRepository{db: db}
}

// GormLogItemRepository reposito logItem
type GormLogItemRepository struct {
	db *gorm.DB
}

// CreateRecord create new record of logItem
func (r GormLogItemRepository) CreateRecord(item Item, MessageID int32, senderID int32) (*LogItem, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	logItem := LogItem{}

	logItem.ID = uid.String()
	logItem.Name = item.Name
	logItem.Amount = item.Amount
	logItem.Category = item.Category
	logItem.MessageID = MessageID
	logItem.CreatedAt = Timestamp()
	logItem.TelegramUserID = senderID

	logrus.Info("Create record logItem ", logItem)
	err = r.db.Create(&logItem).Error
	if err != nil {
		return nil, err
	}
	return &logItem, err
}

// UpdateRecord update record of logItem
func (r GormLogItemRepository) UpdateRecord(logItem *LogItem, item Item, senderID int32) error {
	if logItem.ID == "" {
		return errors.New("can update only created items")
	}

	logItem.Name = item.Name
	logItem.Amount = item.Amount
	logItem.Category = item.Category

	if err := r.db.Save(&logItem).Error; err != nil {
		return err
	}
	logrus.Info("Update record logItem ", logItem)

	return nil
}

// GetRecordsByTelegramID get message by message id
func (r GormLogItemRepository) GetRecordsByTelegramID(SenderID int32) ([]LogItem, error) {
	var items []LogItem
	if err := r.db.Where("telegram_user_id = ?", SenderID).Order("created_at").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetRecordsByTelegramIDCurrentMonth delete's message by message id
func (r GormLogItemRepository) GetRecordsByTelegramIDCurrentMonth(SenderID int32) ([]LogItem, error) {
	var items []LogItem

	beginOfMonth := uint64(now.BeginningOfMonth().UnixNano() / int64(time.Second))

	if err := r.db.Where("telegram_user_id = ? AND created_at >= ?", SenderID, beginOfMonth).Order("created_at").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// DeleteRecordsByMessageID delete's messages by message id
func (r GormLogItemRepository) DeleteRecordsByMessageID(MessageID int32) error {
	return r.db.Where("message_id = ?", MessageID).Delete(&LogItem{}).Error
}

// RecordExists check if message's log item is processed and saved into db
func (r GormLogItemRepository) RecordExists(MessageID int32) bool {
	return !r.db.Where("message_id = ?", MessageID).First(&LogItem{}).RecordNotFound()
}

// FetchMostRelevantCategory get most relevant category using count
func (r GormLogItemRepository) FetchMostRelevantCategory(name string, telegramUserID int32) (string, error) {
	type Result struct {
		Category string
	}
	var result Result

	err := r.db.Raw(`
	SELECT
       log_items.category as category,
       COUNT(*) AS count
	FROM log_items 
	WHERE log_items.name = ? AND log_items.telegram_user_id = ?
	GROUP BY log_items.category
	ORDER BY count DESC 
	LIMIT 1;`, name, telegramUserID).Scan(&result).Error
	if err != nil {
		return "", err
	}
	return result.Category, nil
}
