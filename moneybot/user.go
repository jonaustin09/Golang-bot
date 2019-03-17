package moneybot

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

// User stores telegram user's data in db
type User struct {
	ID           uint64
	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
	CreatedAt    uint64
}

// Recipient to implement interface Recipient that allow to send message to
func (user User) Recipient() string {
	return strconv.Itoa(int(user.ID))
}

func (user *User) fetchByID(uid uint64) error {
	return Db.First(&user, uid).Error
}

func (user *User) fetchOrCreate(uid uint64, firstName string,
	lastName string, languageCode string, username string) bool {
	if Db.First(&user, uid).RecordNotFound() {
		user.ID = uid
		user.FirstName = firstName
		user.LastName = lastName
		user.LanguageCode = languageCode
		user.Username = username
		user.CreatedAt = timestamp()
		Db.Create(&user)
		logrus.Info("Create record user", user)
		return true
	}

	return false
}
