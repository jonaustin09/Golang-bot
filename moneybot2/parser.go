package moneybot2

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
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

func ParsedDataIsValid(parsedData []ParsedData) bool {
	return len(parsedData) > 0
}
