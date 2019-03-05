package main

import (
	"testing"
)

func TestParsedData_HasCategory(t *testing.T) {
	text := "item 1"

	items := GetParsedData(text)

	if items[0].HasCategory() {
		t.Errorf("Should not not have category, got: %v", items)
	}

	text = "item 1 category name"

	items = GetParsedData(text)

	if !items[0].HasCategory() {
		t.Errorf("Should have category, got: %v", items)
	}
	if items[0].Category != "category name" {
		t.Errorf("Category should be: %s, got: %s", "category name", items[0].Category)
	}

	text = "item 1 category name\n item 2 category another name\n item 3"

	items = GetParsedData(text)

	if !items[0].HasCategory() {
		t.Errorf("Should have category, got: %v", items)
	}
	if items[0].Category != "category name" {
		t.Errorf("Category should be: %s, got: %s", "category name", items[0].Category)
	}

	if !items[1].HasCategory() {
		t.Errorf("Should have category, got: %v", items)
	}
	if items[1].Category != "category another name" {
		t.Errorf("Category should be: %s, got: %s", "category another name", items[1].Category)
	}

	if items[2].HasCategory() {
		t.Errorf("Should not have category, got: %v", items)
	}
}

func TestParsedData_IsValid(t *testing.T) {
	text := "item"
	items := GetParsedData(text)
	if len(items) != 0 {
		t.Errorf("Should not be valid, got: %v", items)
	}

	text = "item 1d"
	items = GetParsedData(text)
	if len(items) != 0 {
		t.Errorf("Should not be valid, got: %v", items)
	}

	text = ""
	items = GetParsedData(text)
	if len(items) != 0 {
		t.Errorf("Should not be valid, got: %v", items)
	}

	text = ""
	items = GetParsedData(text)
	if len(items) != 0 {
		t.Errorf("Should not be valid, got: %v", items)
	}

	text = "test 1"
	items = GetParsedData(text)
	if len(items) != 1 {
		t.Errorf("Should be valid, got: %v", items)
	}

	text = "test 1 category"
	items = GetParsedData(text)
	if len(items) != 1 {
		t.Errorf("Should be valid, got: %v", items)
	}

	text = "test with long name 10.1 category is not too short also"
	items = GetParsedData(text)
	if len(items) != 1 {
		t.Errorf("Should be valid, got: %v", items)
	}

	text = "test 1 category\n test2 category\n test 3invalid"
	items = GetParsedData(text)
	if len(items) != 0 {
		t.Errorf("Should not be valid, got: %v", items)
	}
}

func TestParsedData_ParsedDataIsValid(t *testing.T) {
	text := "test 1 category\n test2 category\n test 3invalid\n test 4"
	items := GetParsedData(text)

	if parsedDataIsValid(items) {
		t.Errorf("Should be not valid, got: %v", items)
	}
}
