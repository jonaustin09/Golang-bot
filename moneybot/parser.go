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

// ParsedData contains info about message data
type ParsedData struct {
	Name     string
	Amount   float64
	Category string
}

// IsValid validate parsed data
func (p *ParsedData) IsValid() bool {
	return p.Name != "" && p.Amount != 0
}

// HasCategory get info if repository is set
func (p *ParsedData) HasCategory() bool {
	return p.Category != ""
}

func (p ParsedData) ProcessSaving(messageID int32, sender *tb.User, b *tb.Bot, lr LogItemRepository, config Config) (*LogItem, error) {
	var err error
	if p.Category == "" {
		p.Category, err = lr.FetchMostRelevantCategory(p.Name, int32(sender.ID))
		if err != nil {
			logrus.Error(err)
		}
	}
	logp, err := lr.CreateRecord(p, int32(messageID), int32(sender.ID))
	if err != nil {
		return nil, err
	}

	return logp, nil
}

func SaveParsedData(parsedData []ParsedData, messageID int32, sender *tb.User, b *tb.Bot, lr LogItemRepository, config Config) {
	var sum float64

	var text strings.Builder

	for _, item := range parsedData {
		logp, err := item.ProcessSaving(messageID, sender, b, lr, config)
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

	err := SendServiceMessage(sender, b, text.String(), config.NotificationTimeout)
	if err != nil {
		logrus.Error(err)
	}
}

// GetParsedData parse data from user input
func GetParsedData(s string) []ParsedData {
	var parsedData []ParsedData

	for _, item := range strings.Split(s, "\n") {
		var data ParsedData
		match := myExp.FindStringSubmatch(item)
		if match == nil {
			return []ParsedData{}
		}

		amountStr := strings.Replace(match[amountIndex], ",", ".", 1)
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			logrus.Error(err)
			return []ParsedData{}
		}

		data.Name = strings.TrimSpace(match[nameIndex])
		data.Category = strings.TrimSpace(match[categoryIndex])
		data.Amount = amount
		if data.IsValid() {
			parsedData = append(parsedData, data)
		} else {
			return []ParsedData{}
		}

	}

	return parsedData
}
