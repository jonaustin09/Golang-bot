package moneybot

import (
	"fmt"
	"github.com/dobrovolsky/money_bot/stats"
	"strconv"
)

// User struct for passing to bot
type User struct {
	ID int32
}

// Recipient to implement interface Recipient that allow to send message to
func (user User) Recipient() string {
	return strconv.Itoa(int(user.ID))
}

// LogItem store information about log
type LogItem struct {
	ID        string
	CreatedAt int32
	Name      string
	Amount    float64
	MessageID int32
	Category  string
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

	return fmt.Sprintf("%s %s %.2f %s", timeString, l.Name, l.Amount, l.getCategoryNameOrDefault())
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
	}
}

// PrepareForAnalyze creates message for gRPC
func PrepareForAnalyze(items []LogItem) []*stats.LogMessageAggregated {
	itemsForAnalyze := make([]*stats.LogMessageAggregated, 0, len(items))
	for _, item := range items {
		itemsForAnalyze = append(itemsForAnalyze, &stats.LogMessageAggregated{
			CreatedAt: int64(item.CreatedAt),
			Amount:    float32(item.Amount),
			Category:  item.getCategoryNameOrDefault(),
		})
	}
	return itemsForAnalyze
}
