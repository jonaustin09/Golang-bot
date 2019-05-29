package moneybot2

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/sirupsen/logrus"
)

// Repository represent the category repository contract
type LogItemRepository interface {
	CreateRecord(parsedData ParsedData, MessageID int32, senderID int32) (*LogItem, error)
	UpdateRecord(logItem *LogItem, parsedData ParsedData, senderID int32) error
	DeleteRecordsByMessageID(MessageID int32) error
	GetRecordsByTelegramIDCurrentMonth(SenderID int32) ([]LogItem, error)
	GetRecordsByTelegramID(SenderID int32) ([]LogItem, error)
	RecordExists(MessageID int32) bool
}

func NewGormLogItemRepository(db *gorm.DB) LogItemRepository {
	return GormLogItemRepository{db: db}
}

type GormLogItemRepository struct {
	db *gorm.DB
}

func (r GormLogItemRepository) CreateRecord(parsedData ParsedData, MessageID int32, senderID int32) (*LogItem, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	logItem := LogItem{}

	logItem.ID = uid.String()
	logItem.Name = parsedData.Name
	logItem.Amount = parsedData.Amount
	logItem.Category = parsedData.Category
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

func (r GormLogItemRepository) UpdateRecord(logItem *LogItem, parsedData ParsedData, senderID int32) error {
	if logItem.ID == "" {
		return errors.New("can update only created items")
	}

	logItem.Name = parsedData.Name
	logItem.Amount = parsedData.Amount
	logItem.Category = parsedData.Category

	if err := r.db.Save(&logItem).Error; err != nil {
		return err
	}
	logrus.Info("Update record logItem ", logItem)

	return nil
}

func (r GormLogItemRepository) GetRecordsByTelegramID(SenderID int32) ([]LogItem, error) {
	var items []LogItem
	if err := r.db.Where("telegram_user_id = ?", SenderID).Order("created_at").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r GormLogItemRepository) GetRecordsByTelegramIDCurrentMonth(SenderID int32) ([]LogItem, error) {
	var items []LogItem

	beginOfMonth := uint64(now.BeginningOfMonth().UnixNano() / int64(time.Second))

	if err := r.db.Where("telegram_user_id = ? AND created_at >= ?", SenderID, beginOfMonth).Order("created_at").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r GormLogItemRepository) DeleteRecordsByMessageID(MessageID int32) error {
	return r.db.Where("message_id = ?", MessageID).Delete(&LogItem{}).Error
}

func (r GormLogItemRepository) RecordExists(MessageID int32) bool {
	return !r.db.Where("message_id = ?", MessageID).First(&LogItem{}).RecordNotFound()
}
