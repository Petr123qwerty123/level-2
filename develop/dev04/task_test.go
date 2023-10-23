package main

import (
	"reflect"
	"testing"
)

func TestFindAnagramSets001(t *testing.T) {
	input := []string{"листок", "пятак", "пятка", "слиток", "столик", "тяпка"}
	expected := map[string][]string{
		"листок": {"листок", "слиток", "столик"},
		"пятак":  {"пятак", "пятка", "тяпка"},
	}

	actual := FindAnagramSets(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestFindAnagramSets002(t *testing.T) {
	input := []string{"ЛиСтОк", "ПяТаК", "ПяТкА", "СлИтОк", "СтОлИк", "ТяПкА"}
	expected := map[string][]string{
		"ЛиСтОк": {"листок", "слиток", "столик"},
		"ПяТаК":  {"пятак", "пятка", "тяпка"},
	}

	actual := FindAnagramSets(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestFindAnagramSets003(t *testing.T) {
	input := []string{"листок", "пятак"}
	expected := map[string][]string{}

	actual := FindAnagramSets(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestFindAnagramSets004(t *testing.T) {
	input := []string{
		"ЛиСтОк",
		"лИсТоК",
		"ПяТаК",
		"пЯтАк",
		"ПяТкА",
		"пЯтКа",
		"СлИтОк",
		"сЛиТоК",
		"СтОлИк",
		"сТоЛиК",
		"ТяПкА",
		"тЯпКа",
	}
	expected := map[string][]string{
		"ЛиСтОк": {"листок", "слиток", "столик"},
		"ПяТаК":  {"пятак", "пятка", "тяпка"},
	}

	actual := FindAnagramSets(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}
