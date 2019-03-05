package main

import (
	"github.com/sirupsen/logrus"
)

// Category stores users' category
type Category struct {
	ID             uint64
	Name           string `gorm:"INDEX"`
	TelegramUserID uint64 `gorm:"INDEX"`
	createdAt      uint64
}

func (c *Category) fetchOrCreate(name string, telegramUserID uint64) error {
	if db.Where("name = ? AND telegram_user_id = ?", name, telegramUserID).First(c).RecordNotFound() {
		c.Name = name
		c.TelegramUserID = telegramUserID
		c.createdAt = timestamp()
		logrus.Info("Create new category", c)
		return db.Create(&c).Error
	}
	return nil
}

func (c *Category) fetchByID(ID uint64) error {
	return db.First(c, ID).Error
}

func (c *Category) getDefault() error {
	return db.First(c, 9999).Error
}

func (c *Category) fetchMostRelevantForItem(name string, telegramUserID uint64) error {
	return db.Raw(
		`SELECT id, name, MAX(c) as _count FROM 
				(SELECT categories.id, categories.name , COUNT(log_items.category_id) as c FROM categories 
					JOIN log_items ON log_items.category_id = categories.id
					WHERE log_items.name = ? AND log_items.telegram_user_id = ?
					GROUP BY log_items.category_id);`,
		name, telegramUserID).Scan(&c).Error
}
