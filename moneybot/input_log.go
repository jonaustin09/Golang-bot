package moneybot

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// InputLogRepository represent the inputLog repository contract
type InputLogRepository interface {
	CreateRecord(text string, senderID int32) error
}

// NewGormInputLogRepository creates new repository
func NewGormInputLogRepository(db *gorm.DB) InputLogRepository {
	return GormInputLogRepository{db: db}
}

// GormInputLogRepository is repository for saving imputed text to db
type GormInputLogRepository struct {
	db *gorm.DB
}

// CreateRecord save to db what user entered
func (r GormInputLogRepository) CreateRecord(text string, senderID int32) error {
	var inputLogger InputLog
	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	inputLogger.ID = uid.String()
	inputLogger.Text = text
	inputLogger.TelegramUserID = senderID
	inputLogger.CreatedAt = Timestamp()

	logrus.Info("Create record inputLogger ", inputLogger)
	err = r.db.Create(&inputLogger).Error
	return err
}
