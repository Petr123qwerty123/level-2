package main

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
	"wb-level-2/develop/dev03/utils"
)

func TestIntSlice_UnmarshalText(t *testing.T) {
	t.Run("Valid input", func(t *testing.T) {
		input := []byte("1,2,3")
		expected := IntSlice{1, 2, 3}
		var is IntSlice

		err := is.UnmarshalText(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(is, expected) {
			t.Errorf("got %v, want %v", is, expected)
		}
	})

	t.Run("Invalid input", func(t *testing.T) {
		input := []byte("1,2-")
		var is IntSlice

		err := is.UnmarshalText(input)
		if !errors.Is(err, errParse) {
			t.Errorf("got %v, want %v", err, errParse)
		}
	})
}

// Helper function to reset the command-line args.
func resetArgs() {
	os.Args = []string{"testArgs"}
}

// Helper function to set the command-line args to sf fields
func setAndParseArgs(args []string, sf *SortFlags) {
	for _, arg := range args {
		os.Args = append(os.Args, arg)
	}

	sf.Parse()
}

func TestSortFlags_Parse(t *testing.T) {
	var sf SortFlags

	args := []string{"-k", "1,2", "-n", "-r", "-u", "-M", "-b", "-c", "-h"}

	t.Run("Parse flags", func(t *testing.T) {
		resetArgs()
		setAndParseArgs(args, &sf)

		expectedColumns := IntSlice{1, 2}
		expectedNumeric := true
		expectedReverse := true
		expectedUnique := true
		expectedMonth := true
		expectedIgnoreSpaces := true
		expectedCheckSorted := true
		expectedNumericSuffix := true

		if !reflect.DeepEqual(sf.columns, expectedColumns) {
			t.Errorf("got %v, want %v", sf.columns, expectedColumns)
		}

		if sf.numeric != expectedNumeric {
			t.Errorf("got %t, want %t", sf.numeric, expectedNumeric)
		}

		if sf.reverse != expectedReverse {
			t.Errorf("got %t, want %t", sf.reverse, expectedReverse)
		}

		if sf.unique != expectedUnique {
			t.Errorf("got %t, want %t", sf.unique, expectedUnique)
		}

		if sf.month != expectedMonth {
			t.Errorf("got %t, want %t", sf.month, expectedMonth)
		}

		if sf.ignoreSpaces != expectedIgnoreSpaces {
			t.Errorf("got %t, want %t", sf.ignoreSpaces, expectedIgnoreSpaces)
		}

		if sf.checkSorted != expectedCheckSorted {
			t.Errorf("got %t, want %t", sf.checkSorted, expectedCheckSorted)
		}

		if sf.numericSuffix != expectedNumericSuffix {
			t.Errorf("got %t, want %t", sf.numericSuffix, expectedNumericSuffix)
		}
	})
}

func TestSortClient_Sort(t *testing.T) {
	t.Run("Sort data", func(t *testing.T) {
		sc := &SortClient{
			flags: SortFlags{
				columns:       IntSlice{1},
				numeric:       true,
				reverse:       false,
				unique:        false,
				month:         false,
				ignoreSpaces:  false,
				checkSorted:   false,
				numericSuffix: false,
			},
			data: []string{"3 4", "1 5", "2 3"},
		}

		expected := []string{"1 5", "2 3", "3 4"}

		result, err := sc.Sort()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, want %v", result, expected)
		}
	})
}

func TestSortClient_IsSorted(t *testing.T) {
	t.Run("Data sorted", func(t *testing.T) {
		sc := &SortClient{
			flags: SortFlags{},
			data:  []string{"apple", "banana", "cherry"},
		}

		index, err := sc.IsSorted()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if index != -1 {
			t.Errorf("expected -1, got %d", index)
		}
	})

	t.Run("Data not sorted", func(t *testing.T) {
		sc := &SortClient{
			flags: SortFlags{},
			data:  []string{"banana", "apple", "cherry"},
		}

		expectedIndex := 0

		index, err := sc.IsSorted()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if index != expectedIndex {
			t.Errorf("expected %d, got %d", expectedIndex, index)
		}
	})
}

func TestReadData(t *testing.T) {
	t.Run("Read data", func(t *testing.T) {
		input := strings.NewReader("line1\nline2\nline3")
		expected := []string{"line1", "line2", "line3"}

		result, err := utils.ReadData(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, want %v", result, expected)
		}
	})
}

func TestWriteData(t *testing.T) {
	t.Run("Write data", func(t *testing.T) {
		var output strings.Builder
		data := []string{"line1", "line2", "line3"}

		err := utils.WriteData(&output, data...)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := "line1\nline2\nline3\n"

		if output.String() != expected {
			t.Errorf("got %q, want %q", output.String(), expected)
		}
	})
}

func TestTrimLeadingSpaces(t *testing.T) {
	t.Run("Trim leading spaces", func(t *testing.T) {
		input := "   \n\ttrim spaces  "
		expected := "trim spaces  "

		result := utils.TrimLeadingSpaces(input)
		if result != expected {
			t.Errorf("got %q, want %q", result, expected)
		}
	})
}

func TestParseMonth(t *testing.T) {
	t.Run("Valid month", func(t *testing.T) {
		input := "   Jan   "
		expected := 0

		result, err := utils.ParseMonth(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if result != expected {
			t.Errorf("got %d, want %d", result, expected)
		}
	})

	t.Run("Invalid month", func(t *testing.T) {
		input := "   XYZ   "

		expectedError := errors.New("month cannot be parsed")

		_, err := utils.ParseMonth(input)
		if err.Error() != expectedError.Error() {
			t.Errorf("got %v, want %v", err, expectedError)
		}
	})
}

func TestParseNumericValue(t *testing.T) {
	t.Run("With suffix", func(t *testing.T) {
		input := "10K"
		suffix := true
		expected := 10000.0

		result, err := utils.ParseNumericValue(input, suffix)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if result != expected {
			t.Errorf("got %f, want %f", result, expected)
		}
	})

	t.Run("Without suffix", func(t *testing.T) {
		input := "10"
		suffix := false
		expected := 10.0

		result, err := utils.ParseNumericValue(input, suffix)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if result != expected {
			t.Errorf("got %f, want %f", result, expected)
		}
	})

	t.Run("Invalid value", func(t *testing.T) {
		input := "abc"
		suffix := false

		expectedError := errors.New("number cannot be parsed")

		_, err := utils.ParseNumericValue(input, suffix)
		if err.Error() != expectedError.Error() {
			t.Errorf("got %v, want %v", err, expectedError)
		}
	})
}

func TestRemoveDuplicates(t *testing.T) {
	t.Run("Remove duplicates", func(t *testing.T) {
		input := []string{"line1", "line2", "line1", "line3"}
		expected := []string{"line1", "line2", "line3"}

		result := utils.RemoveDuplicates(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, want %v", result, expected)
		}
	})
}
