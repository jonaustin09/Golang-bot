package moneybot2

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type InputLogRepository interface {
	CreateRecord(text string, senderID int32) error
}

func NewGormInputLogRepository(db *gorm.DB) InputLogRepository {
	return GormInputLogRepository{db: db}
}

type GormInputLogRepository struct {
	db *gorm.DB
}

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
