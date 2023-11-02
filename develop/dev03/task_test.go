package main

import (
	"errors"
	"flag"
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
func resetArgs(args []string) {
	os.Args = []string{"testArgs"}

	for _, arg := range args {
		os.Args = append(os.Args, arg)
	}
}

func TestSortFlags_Parse(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		flags SortFlags
	}{
		{
			name:  "No flags",
			args:  []string{},
			flags: SortFlags{},
		},
		{
			name: "Columns flag",
			args: []string{"-k", "1,3,5-7,9-13,15"},
			flags: SortFlags{
				columns: IntSlice{1, 3, 5, 6, 7, 9, 10, 11, 12, 13, 15},
			},
		},
		{
			name: "Numeric flag",
			args: []string{"-n"},
			flags: SortFlags{
				numeric: true,
			},
		},
		{
			name: "Reverse flag",
			args: []string{"-r"},
			flags: SortFlags{
				reverse: true,
			},
		},
		{
			name: "Unique flag",
			args: []string{"-u"},
			flags: SortFlags{
				unique: true,
			},
		},
		{
			name: "Month flag",
			args: []string{"-M"},
			flags: SortFlags{
				month: true,
			},
		},
		{
			name: "Ignore spaces flag",
			args: []string{"-b"},
			flags: SortFlags{
				ignoreSpaces: true,
			},
		},
		{
			name: "Check sorted flag",
			args: []string{"-c"},
			flags: SortFlags{
				checkSorted: true,
			},
		},
		{
			name: "Numeric suffix flag",
			args: []string{"-h"},
			flags: SortFlags{
				numericSuffix: true,
			},
		},
		{
			name:  "Invalid flag",
			args:  []string{"-x"},
			flags: SortFlags{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			sf := SortFlags{}

			resetArgs(tt.args)
			sf.Parse()

			if !reflect.DeepEqual(sf, tt.flags) {
				t.Errorf("SortFlags.Parse() got = %v, want %v", sf, tt.flags)
			}
		})
	}
}

func TestSortArgs_Parse(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		sortArgs SortArgs
		hasErr   bool
	}{
		{
			name: "No args",
			args: []string{},
			sortArgs: SortArgs{
				inputFiles: []*os.File{os.Stdin},
			},
			hasErr: false,
		},
		{
			name:     "Invalid filename",
			args:     []string{"kjhg??fd"},
			sortArgs: SortArgs{},
			hasErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			sa := SortArgs{}

			resetArgs(tt.args)

			flag.Parse()
			err := sa.Parse()

			if (err != nil) != tt.hasErr {
				t.Errorf("SortArgs.Parse() error = %v, wantErr %v", err, tt.hasErr)
				return
			}

			if len(sa.inputFiles) != len(tt.sortArgs.inputFiles) {
				t.Errorf("SortArgs.Parse() inputFiles length = %v, want %v", len(sa.inputFiles), len(tt.sortArgs.inputFiles))
			}
		})
	}
}

func TestSortClient_Sort(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		data           []string
		expectedResult []string
	}{
		{
			name:           "Sort with no flags",
			args:           []string{},
			data:           []string{"b", "c", "a", "e", "f", "d"},
			expectedResult: []string{"a", "b", "c", "d", "e", "f"},
		},
		{
			name: "Sort with columns flag",
			args: []string{"-k", "2,4"},
			data: []string{
				"2 2 1 1",
				"1 2 2 2",
				"2 2 1 2",
				"2 2 2 1",
				"1 2 1 2",
				"1 2 2 1",
				"1 2 1 1",
				"2 1 1 1",
				"2 1 2 2",
				"2 2 2 2",
				"1 1 1 1",
				"1 1 1 2",
				"1 1 2 1",
				"2 1 1 2",
				"2 1 2 1",
				"1 1 2 2",
			},
			expectedResult: []string{
				"2 1 1 1",
				"2 1 2 1",
				"1 1 2 1",
				"1 1 1 1",
				"2 1 2 2",
				"1 1 1 2",
				"2 1 1 2",
				"1 1 2 2",
				"1 2 2 1",
				"1 2 1 1",
				"2 2 1 1",
				"2 2 2 1",
				"1 2 1 2",
				"2 2 2 2",
				"2 2 1 2",
				"1 2 2 2",
			},
		},
		{
			name: "Sort with numeric flag",
			args: []string{"-n"},
			data: []string{
				"1",
				"20",
				"22",
				"200",
				"100",
				"10",
				"11",
				"2",
			},
			expectedResult: []string{
				"1",
				"2",
				"10",
				"11",
				"20",
				"22",
				"100",
				"200",
			},
		},
		{
			name: "Sort with reverse flag",
			args: []string{"-r"},
			data: []string{
				"1",
				"20",
				"22",
				"200",
				"100",
				"10",
				"11",
				"2",
			},
			expectedResult: []string{
				"22",
				"200",
				"20",
				"2",
				"11",
				"100",
				"10",
				"1",
			},
		},
		{
			name: "Sort with unique flag",
			args: []string{"-u"},
			data: []string{
				"1",
				"20",
				"22",
				"200",
				"100",
				"10",
				"22",
				"200",
				"100",
				"11",
				"2",
				"1",
				"20",
				"200",
				"100",
				"10",
				"22",
				"200",
				"100",
				"11",
				"2",
				"22",
				"100",
				"10",
				"22",
				"200",
				"100",
				"11",
				"2",
				"1",
				"20",
				"200",
				"100",
				"10",
			},
			expectedResult: []string{
				"1",
				"10",
				"100",
				"11",
				"2",
				"20",
				"200",
				"22",
			},
		},
		{
			name: "Sort with month flag",
			args: []string{"-M"},
			data: []string{
				"mArcH",
				"AprIl",
				"JaNuary",
				"fEbRuary",
				"October",
				"nOVember",
				"DecEMber",
				"mAY",
				"July",
				"jUNe",
				"SepTEMber",
				"auGUst",
			},
			expectedResult: []string{
				"JaNuary",
				"fEbRuary",
				"mArcH",
				"AprIl",
				"mAY",
				"jUNe",
				"July",
				"auGUst",
				"SepTEMber",
				"October",
				"nOVember",
				"DecEMber",
			},
		},
		{
			name: "Sort with ignore spaces flag",
			args: []string{"-b"},
			data: []string{
				"					b",
				"		c",
				"				d",
				"						e",
				"	f",
				"		a",
				"			g",
				"							k",
				"	h",
				"	i",
				"			j",
				"l",
			},
			expectedResult: []string{
				"		a",
				"					b",
				"		c",
				"				d",
				"						e",
				"	f",
				"			g",
				"	h",
				"	i",
				"			j",
				"							k",
				"l",
			},
		},
		{
			name: "Sort with numeric suffix flag",
			args: []string{"-h"},
			data: []string{
				"1T",
				"2G",
				"3M",
				"4K",
				"5",
				"6",
			},
			expectedResult: []string{
				"5",
				"6",
				"4K",
				"3M",
				"2G",
				"1T",
			},
		},
		{
			name: "Sort with columns and numeric flag",
			args: []string{"-k", "2", "-n"},
			data: []string{
				"b 5",
				"d 6",
				"h 7",
				"e 8",
				"a 4",
				"c 3",
				"f 2",
				"g 1",
			},
			expectedResult: []string{
				"g 1",
				"f 2",
				"c 3",
				"a 4",
				"b 5",
				"d 6",
				"h 7",
				"e 8",
			},
		},
		{
			name: "Sort with columns, numeric and reverse flag",
			args: []string{"-k", "2", "-n", "-r"},
			data: []string{
				"b 5",
				"d 6",
				"b 1",
				"d 2",
				"e 5",
				"g 8",
				"h 7",
				"e 8",
				"a 1",
				"a 4",
				"c 3",
				"f 2",
				"g 1",
			},
			expectedResult: []string{
				"g 8",
				"e 8",
				"h 7",
				"d 6",
				"e 5",
				"b 5",
				"a 4",
				"c 3",
				"f 2",
				"d 2",
				"b 1",
				"g 1",
				"a 1",
			},
		},
		{
			name: "Sort with columns and month flag",
			args: []string{"-k", "2", "-M"},
			data: []string{
				"b may",
				"d june",
				"d february",
				"h july",
				"e august",
				"a january",
				"a april",
				"c march",
			},
			expectedResult: []string{
				"a january",
				"d february",
				"c march",
				"a april",
				"b may",
				"d june",
				"h july",
				"e august",
			},
		},
		{
			name: "Sort with columns, month and reverse flag",
			args: []string{"-k", "2", "-M", "-r"},
			data: []string{
				"b may",
				"d june",
				"d february",
				"h july",
				"e august",
				"a january",
				"a april",
				"c march",
			},
			expectedResult: []string{
				"e august",
				"h july",
				"d june",
				"b may",
				"a april",
				"c march",
				"d february",
				"a january",
			},
		},
		{
			name: "Sort with columns and numeric suffix flag",
			args: []string{"-k", "2", "-h"},
			data: []string{
				"b 5g",
				"d 6t",
				"h 7o",
				"e 8k",
				"a 4g",
				"c 3m",
				"f 2k",
				"g 1t",
			},
			expectedResult: []string{
				"h 7o",
				"f 2k",
				"e 8k",
				"c 3m",
				"a 4g",
				"b 5g",
				"g 1t",
				"d 6t",
			},
		},
		{
			name: "Sort with columns, numeric suffix and reverse flag",
			args: []string{"-k", "2", "-h", "-r"},
			data: []string{
				"b 5g",
				"d 6t",
				"h 7o",
				"e 8k",
				"a 4g",
				"c 3m",
				"f 2k",
				"g 1t",
			},
			expectedResult: []string{
				"d 6t",
				"g 1t",
				"b 5g",
				"a 4g",
				"c 3m",
				"e 8k",
				"f 2k",
				"h 7o",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			resetArgs(tt.args)

			sc := SortClient{}

			sc.flags.Parse()
			err := sc.args.Parse()
			if err != nil {
				t.Errorf("not expected error: %q", err)
			}

			sc.data = tt.data

			result := sc.Sort()

			if len(result) == 0 && len(tt.expectedResult) == 0 {
				return
			}

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("got %v, expected %v", result, tt.expectedResult)
			}
		})
	}
}

func TestSortClient_Sort2(t *testing.T) {
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

		result := sc.Sort()

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

		index := sc.IsSorted()

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

		index := sc.IsSorted()

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
