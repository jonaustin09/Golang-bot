package moneybot

import (
	"fmt"
)

// LogItem store information about log
type LogItem struct {
	ID             string
	CreatedAt      int32
	Name           string
	Amount         float64
	MessageID      int32
	TelegramUserID string
	Category       string
}

func (l *LogItem) getCategoryNameOrDefault() string {
	if l.Category == "" {
		return "other"
	}
	return l.Category
}

// String string presentation
func (l *LogItem) String() string {
	localTime := GetLocalTime(l.CreatedAt)
	timeString := localTime.Format("02.01.2006")

	return fmt.Sprintf("%s %s %.2f %s %s", timeString, l.Name, l.Amount, l.getCategoryNameOrDefault(), l.TelegramUserID)
}

// toCSV gets csv data
func (l *LogItem) toCSV() []string {
	localTime := GetLocalTime(l.CreatedAt)
	timeString := localTime.Format("02.01.2006")

	return []string{
		timeString,
		l.Name,
		fmt.Sprintf("%.2f", l.Amount),
		l.getCategoryNameOrDefault(),
		l.TelegramUserID,
	}
}
