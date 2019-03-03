package main

// Category stores users' category
type Category struct {
	ID             uint64
	Name           string `gorm:"INDEX"`
	TelegramUserID uint64 `gorm:"INDEX"`
	CreatedAt      uint64
}

func (c *Category) fetchOrCreate(name string, telegramUserID uint64) error {
	if db.Where("name = ? AND telegram_user_id = ?", name, telegramUserID).First(c).RecordNotFound() {
		c.Name = name
		c.TelegramUserID = telegramUserID
		c.CreatedAt = Timestamp()
		return db.Create(c).Error
	}
	return nil
}

func (c *Category) fetchByID(ID uint64) error {
	return db.First(c, ID).Error
}
