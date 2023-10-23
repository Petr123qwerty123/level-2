package main

import (
	"errors"
	"testing"
)

func TestUnpack001(t *testing.T) {
	input := "a4bc2d5e"
	expected := "aaaabccddddde"

	actual, err := Unpack(input)

	if err != nil {
		t.Errorf("Should not produce an error")
	}

	if expected != actual {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestUnpack002(t *testing.T) {
	input := "abcd"
	expected := "abcd"

	actual, err := Unpack(input)

	if err != nil {
		t.Error("Should not produce an error")
	}

	if expected != actual {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestUnpack003(t *testing.T) {
	input := "45"
	expectedError := errors.New("invalid string")

	_, err := Unpack(input)

	if err.Error() != expectedError.Error() {
		t.Errorf("Expected %q, got %q", expectedError, err)
	}
}

func TestUnpack004(t *testing.T) {
	input := ""
	expected := ""

	actual, err := Unpack(input)

	if err != nil {
		t.Error("Should not produce an error")
	}

	if expected != actual {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestUnpack005(t *testing.T) {
	input := "qwe\\4\\5"
	expected := "qwe45"

	actual, err := Unpack(input)

	if err != nil {
		t.Error("Should not produce an error")
	}

	if expected != actual {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestUnpack006(t *testing.T) {
	input := "qwe\\45"
	expected := "qwe44444"

	actual, err := Unpack(input)

	if err != nil {
		t.Error("Should not produce an error")
	}

	if expected != actual {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestUnpack007(t *testing.T) {
	input := "qwe\\\\5"
	expected := "qwe\\\\\\\\\\"

	actual, err := Unpack(input)

	if err != nil {
		t.Error("Should not produce an error")
	}

	if expected != actual {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestUnpack008(t *testing.T) {
	input := `\qwe`
	expectedError := errors.New("invalid string")

	_, err := Unpack(input)

	if err.Error() != expectedError.Error() {
		t.Errorf("Expected %q, got %q", expectedError, err)
	}
}
