package moneybot

import (
	"testing"
)

func TestParsedData_HasCategory(t *testing.T) {
	text := "item 1"

	items := GetParsedData(text)

	if items[0].HasCategory() {
		t.Errorf("Should not not have repository, got: %v", items)
	}

	text = "item 1 repository name"

	items = GetParsedData(text)

	if !items[0].HasCategory() {
		t.Errorf("Should have repository, got: %v", items)
	}
	if items[0].Category != "repository name" {
		t.Errorf("Category should be: %s, got: %s", "repository name", items[0].Category)
	}

	text = "item 1 repository name\n item 2 repository another name\n item 3"

	items = GetParsedData(text)

	if !items[0].HasCategory() {
		t.Errorf("Should have repository, got: %v", items)
	}
	if items[0].Category != "repository name" {
		t.Errorf("Category should be: %s, got: %s", "repository name", items[0].Category)
	}

	if !items[1].HasCategory() {
		t.Errorf("Should have repository, got: %v", items)
	}
	if items[1].Category != "repository another name" {
		t.Errorf("Category should be: %s, got: %s", "repository another name", items[1].Category)
	}

	if items[2].HasCategory() {
		t.Errorf("Should not have repository, got: %v", items)
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

	text = "test 1 repository"
	items = GetParsedData(text)
	if len(items) != 1 {
		t.Errorf("Should be valid, got: %v", items)
	}

	text = "test with long name 10.1 repository is not too short also"
	items = GetParsedData(text)
	if len(items) != 1 {
		t.Errorf("Should be valid, got: %v", items)
	}

	text = "test 1 repository\n test2 repository\n test 3invalid"
	items = GetParsedData(text)
	if len(items) != 0 {
		t.Errorf("Should not be valid, got: %v", items)
	}
}
