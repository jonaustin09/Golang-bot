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

var regex = `(?P<name>[А-ЯҐЄІЇa-яґєшії\s+]+)(?P<amount>\d+((\.|,)\d*)?)\s*(?P<category>[А-ЯҐЄІЇa-яґєшії]{0,})`
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
func GetParsedData(s string) ParsedData {

	parsedData := ParsedData{}

	match := myExp.FindStringSubmatch(s)
	if match == nil {
		return parsedData
	}

	amount_str := strings.Replace(match[amountIndex], ",", ".", 1)
	amount, err := strconv.ParseFloat(amount_str, 64)
	Check(err)

	parsedData.Name = match[nameIndex]
	parsedData.Category = match[categoryIndex]
	parsedData.Amount = amount

	return parsedData
}
