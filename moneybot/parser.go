package moneybot

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

type myRegexp struct {
	*regexp.Regexp
}

const (
	nameIndex = iota + 1
	amountIndex
	categoryIndex = 5
)

var regex = `(?P<name>[А-ЯҐЄІЇa-яґєшії\s\d]+)\s+(?P<amount>\d+((\.|,)\d*)?)(?P<repository>\s+[А-ЯҐЄІЇa-яґєшії\s\d]{0,}|$)`
var myExp = myRegexp{regexp.MustCompile(regex)}

// Item contains info about message data
type Item struct {
	Name     string
	Amount   float64
	Category string
}

// IsValid validate parsed data
func (p Item) IsValid() bool {
	return p.Name != "" && p.Amount > 0
}

// HasCategory get info if repository is set
func (p *Item) HasCategory() bool {
	return p.Category != ""
}

func (p Item) ProcessSaving(messageID int32, sender int32, b *tb.Bot, lr LogItemRepository, config Config) (*LogItem, error) {
	var err error
	if p.Category == "" {
		p.Category, err = lr.FetchMostRelevantCategory(p.Name, sender)
		if err != nil {
			logrus.Error(err)
		}
	}
	logp, err := lr.CreateRecord(p, messageID, sender)
	if err != nil {
		return nil, err
	}

	return logp, nil
}

func SaveItems(items []Item, messageID int32, sender *tb.User, b *tb.Bot, lr LogItemRepository, config Config) {
	var sum float64

	var text strings.Builder

	for _, item := range items {
		logp, err := item.ProcessSaving(messageID, int32(sender.ID), b, lr, config)
		if err != nil {
			logrus.Error(err)
		} else {
			sum += item.Amount
		}
		text.WriteString(fmt.Sprintf("`Create: %s`\n", logp.String()))

	}
	if sum > 0 {
		text.WriteString(fmt.Sprintf("`Sum: %v`", sum))
	}

	err := SendDeletableMessage(sender, b, text.String(), config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
	}
}

// GetParsedData parse data from user input
func GetItem(s string) []Item {
	var items []Item

	for _, substring := range strings.Split(s, "\n") {
		var data Item
		match := myExp.FindStringSubmatch(substring)
		if match == nil {
			return []Item{}
		}

		amountStr := strings.Replace(match[amountIndex], ",", ".", 1)
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			logrus.Error(err)
			return []Item{}
		}

		data.Name = strings.TrimSpace(match[nameIndex])
		data.Category = strings.TrimSpace(match[categoryIndex])
		data.Amount = amount
		if data.IsValid() {
			items = append(items, data)
		} else {
			return []Item{}
		}

	}

	return items
}
