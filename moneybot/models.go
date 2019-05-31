package moneybot

import (
	"fmt"
	"strconv"

	"github.com/dobrovolsky/money_bot/stats"
)

// User stores telegram user's data in db
type User struct {
	ID           int32
	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
	CreatedAt    int32
}

// Recipient to implement interface Recipient that allow to send message to
func (user User) Recipient() string {
	return strconv.Itoa(int(user.ID))
}

// InputLog store input data for future analyze
type InputLog struct {
	ID             string
	Text           string
	TelegramUserID int32
	CreatedAt      int32
}

// LogItem store information about log
type LogItem struct {
	ID             string
	CreatedAt      int32
	Name           string
	Amount         float64
	MessageID      int32
	TelegramUserID int32
	Category       string
}

// String string presentation
func (l *LogItem) String() string {
	localTime := GetLocalTime(l.CreatedAt)
	timeString := localTime.Format("02.01.2006")

	return fmt.Sprintf("%s %s %.2f %s", timeString, l.Name, l.Amount, l.Category)
}

// toCSV get csv data
func (l *LogItem) toCSV() []string {
	localTime := GetLocalTime(l.CreatedAt)
	timeString := localTime.Format("02.01.2006")

	return []string{
		timeString,
		l.Name,
		fmt.Sprintf("%.2f", l.Amount),
		l.Category,
	}
}

// PrepareForAnalyze creates message for gRPC
func PrepareForAnalyze(items []LogItem) []*stats.LogItemMessage {
	itemsForAnalyze := make([]*stats.LogItemMessage, 0, len(items))
	for _, item := range items {
		itemsForAnalyze = append(itemsForAnalyze, &stats.LogItemMessage{
			CreatedAt: int64(item.CreatedAt),
			Name:      item.Name,
			Amount:    float32(item.Amount),
			Category:  item.Category,
		})
	}
	return itemsForAnalyze
}
