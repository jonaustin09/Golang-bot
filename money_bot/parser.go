package money_bot

import (
	"regexp"
	"strconv"
	"strings"
)

type myRegexp struct {
	*regexp.Regexp
}

const (
	nameIndex = iota + 1
	amountIndex
	categoryIndex = 5
)

var regex = `(?P<name>[А-ЯҐЄІЇa-яґєшії\s\d]+)\s+(?P<amount>\d+((\.|,)\d*)?)(?P<category>\s+[А-ЯҐЄІЇa-яґєшії\s\d]{0,}|$)`
var myExp = myRegexp{regexp.MustCompile(regex)}

// ParsedData contains info about message data
type ParsedData struct {
	Name     string
	Amount   float64
	Category string
}

// IsValid validate parsed data
func (p *ParsedData) isValid() bool {
	return p.Name != "" && p.Amount != 0
}

// HasCategory get info if category is set
func (p *ParsedData) hasCategory() bool {
	return p.Category != ""
}

// GetParsedData parse data from user input
func getParsedData(s string) []ParsedData {
	var parsedData []ParsedData

	for _, item := range strings.Split(s, "\n") {
		var data ParsedData
		match := myExp.FindStringSubmatch(item)
		if match == nil {
			return []ParsedData{}
		}

		amountStr := strings.Replace(match[amountIndex], ",", ".", 1)
		amount, err := strconv.ParseFloat(amountStr, 64)
		Check(err)

		data.Name = strings.TrimSpace(match[nameIndex])
		data.Category = strings.TrimSpace(match[categoryIndex])
		data.Amount = amount
		if data.isValid() {
			parsedData = append(parsedData, data)
		} else {
			return []ParsedData{}
		}

	}

	return parsedData
}

func parsedDataIsValid(parsedData []ParsedData) bool {
	return len(parsedData) > 0
}
