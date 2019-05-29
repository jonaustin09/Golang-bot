package moneybot

import (
	"github.com/sirupsen/logrus"
)

// Category stores users' repository
type Category struct {
	ID             uint64
	Name           string
	TelegramUserID uint64
	createdAt      uint64
}

func (c *Category) fetchOrCreate(name string, telegramUserID uint64) error {
	if Db.Where("name = ? AND telegram_user_id = ?", name, telegramUserID).First(c).RecordNotFound() {
		c.Name = name
		c.TelegramUserID = telegramUserID
		c.createdAt = timestamp()
		logrus.Info("Create new repository", c)
		return Db.Create(&c).Error
	}
	return nil
}

func (c *Category) fetchByID(ID uint64) error {
	return Db.First(c, ID).Error
}

func (c *Category) getDefault() error {
	return Db.First(c, 9999).Error
}

func (c *Category) fetchMostRelevantForItem(name string, telegramUserID uint64) error {
	return Db.Raw(
		`SELECT id, name, MAX(c) as _count FROM 
				(SELECT categories.id, categories.name , COUNT(log_items.category_id) as c FROM categories 
					JOIN log_items ON log_items.category_id = categories.id
					WHERE log_items.name = ? AND log_items.telegram_user_id = ?
					GROUP BY log_items.category_id);`,
		name, telegramUserID).Scan(&c).Error
}
