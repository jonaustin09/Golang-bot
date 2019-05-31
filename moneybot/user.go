package moneybot

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// UserRepository represent the user repository contract
type UserRepository interface {
	FetchByID(ID int32) (*User, error)
	FetchOrCreate(ID int32, firstName string,
		lastName string, languageCode string, username string) (*User, error)
}

// NewGormUserRepository create new repository
func NewGormUserRepository(db *gorm.DB) UserRepository {
	return GormUserRepository{db: db}
}

// GormUserRepository repository for saving user into db
type GormUserRepository struct {
	db *gorm.DB
}

// FetchByID get user by id
func (r GormUserRepository) FetchByID(ID int32) (*User, error) {
	var u *User
	err := r.db.First(u, ID).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

// FetchOrCreate create new user if necessary
func (r GormUserRepository) FetchOrCreate(ID int32, firstName string,
	lastName string, languageCode string, username string) (*User, error) {
	var u *User

	if r.db.First(u, ID).RecordNotFound() {
		u := User{}
		u.ID = ID
		u.FirstName = firstName
		u.LastName = lastName
		u.LanguageCode = languageCode
		u.Username = username
		u.CreatedAt = Timestamp()
		err := r.db.Create(&u).Error
		if err != nil {
			return nil, err
		}
		logrus.Info("Create record user", u)
	}

	return u, nil
}
