package moneybot

import (
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

// InputLog store input data for future analyze
type InputLog struct {
	ID             string
	Text           string
	TelegramUserID uint64
	CreatedAt      uint64
}

func (inputLogger *InputLog) createRecord(text string, senderID uint64) error {
	uid, err := uuid.NewV4()
	Check(err)

	inputLogger.ID = uid.String()
	inputLogger.Text = text
	inputLogger.TelegramUserID = senderID
	inputLogger.CreatedAt = timestamp()

	logrus.Info("Create record inputLogger ", inputLogger)
	err = Db.Create(&inputLogger).Error
	return err
}
