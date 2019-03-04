package main

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

var regex = `(?P<name>[А-ЯҐЄІЇa-яґєшії\s+]+)(?P<amount>\d+((\.|,)\d*)?)\s*(?P<category>[А-ЯҐЄІЇa-яґєшії\s+]{0,})`
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

// HasCategory get info if category is set
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
			return parsedData
		}

		amountStr := strings.Replace(match[amountIndex], ",", ".", 1)
		amount, err := strconv.ParseFloat(amountStr, 64)
		Check(err)

		data.Name = match[nameIndex]
		data.Category = match[categoryIndex]
		data.Amount = amount
		parsedData = append(parsedData, data)
	}

	return parsedData
}

func parsedDataIsValid(parsedData []ParsedData) bool {
	isValid := false
	for i := range parsedData {
		isValid = parsedData[i].IsValid()
		if !isValid {
			return false
		}
		isValid = true
	}
	return isValid
}
