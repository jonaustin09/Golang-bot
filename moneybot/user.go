package moneybot

import (
	"strconv"
)

// User struct for passing to bot
type User struct {
	ID int
	Username string
}

// Recipient to implement interface Recipient that allow to send message to
func (user User) Recipient() string {
	return strconv.Itoa(user.ID)
}
