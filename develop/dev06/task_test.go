package main

import (
	"errors"
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"
	"wb-level-2/develop/dev06/utils"
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

func TestCutFlags_Parse(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		flags CutFlags
	}{
		{
			name: "No flags",
			args: []string{},
			flags: CutFlags{
				delimiter: "\t",
			},
		},
		{
			name: "Fields flag",
			args: []string{"-f", "1,3,5-7,9-13,15"},
			flags: CutFlags{
				delimiter: "\t",
				fields:    IntSlice{1, 3, 5, 6, 7, 9, 10, 11, 12, 13, 15},
			},
		},
		{
			name: "Delimiter flag",
			args: []string{"-d", "  "},
			flags: CutFlags{
				delimiter: "  ",
			},
		},
		{
			name: "Separated flag",
			args: []string{"-s"},
			flags: CutFlags{
				delimiter: "\t",
				separated: true,
			},
		},
		{
			name: "Invalid flag",
			args: []string{"-x"},
			flags: CutFlags{
				delimiter: "\t",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			cf := CutFlags{}

			resetArgs(tt.args)
			cf.Parse()

			if !reflect.DeepEqual(cf, tt.flags) {
				t.Errorf("CutFlags.Parse() got = %v, want %v", cf, tt.flags)
			}
		})
	}
}

func TestCutArgs_Parse(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		cutArgs CutArgs
		hasErr  bool
	}{
		{
			name: "No args",
			args: []string{},
			cutArgs: CutArgs{
				inputFiles: []*os.File{os.Stdin},
			},
			hasErr: false,
		},
		{
			name:    "Invalid filename",
			args:    []string{"kjhg??fd"},
			cutArgs: CutArgs{},
			hasErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			ca := CutArgs{}

			resetArgs(tt.args)

			flag.Parse()
			err := ca.Parse()

			if (err != nil) != tt.hasErr {
				t.Errorf("CutArgs.Parse() error = %v, wantErr %v", err, tt.hasErr)
				return
			}

			if len(ca.inputFiles) != len(tt.cutArgs.inputFiles) {
				t.Errorf("CutArgs.Parse() inputFiles length = %v, want %v", len(ca.inputFiles), len(tt.cutArgs.inputFiles))
			}
		})
	}
}

func TestCutClient_Cut(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		data           []string
		expectedResult []string
	}{
		{
			name: "Cut with no flags",
			args: []string{},
			data: []string{
				"b 5 g",
				"d 6 t",
				"h 7 o",
				"e 8 k",
				"a 4 g",
				"c 3 m",
				"f 2 k",
				"g 1 t",
			},
			expectedResult: []string{
				"b 5 g",
				"d 6 t",
				"h 7 o",
				"e 8 k",
				"a 4 g",
				"c 3 m",
				"f 2 k",
				"g 1 t",
			},
		},
		{
			name: "Cut with fields flag",
			args: []string{"-f", "2"},
			data: []string{
				"b\t5\tg",
				"d\t6\tt",
				"h\t7\to",
				"e\t8\tk",
				"a\t4\tg",
				"c\t3\tm",
				"f\t2\tk",
				"g\t1\tt",
			},
			expectedResult: []string{
				"5",
				"6",
				"7",
				"8",
				"4",
				"3",
				"2",
				"1",
			},
		},
		{
			name: "Cut with fields and delimiter flag",
			args: []string{"-f", "2", "-d", " "},
			data: []string{
				"b 5 g",
				"d 6 t",
				"h 7 o",
				"e 8 k",
				"a 4 g",
				"c 3 m",
				"f 2 k",
				"g 1 t",
			},
			expectedResult: []string{
				"5",
				"6",
				"7",
				"8",
				"4",
				"3",
				"2",
				"1",
			},
		},
		{
			name: "Cut with fields and separated flag",
			args: []string{"-f", "2", "-s"},
			data: []string{
				"b5g",
				"d\t6\tt",
				"h7o",
				"e\t8\tk",
				"a4g",
				"c\t3\tm",
				"f2k",
				"g\t1\tt",
			},
			expectedResult: []string{
				"6",
				"8",
				"3",
				"1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			resetArgs(tt.args)

			cc := CutClient{}

			cc.flags.Parse()
			err := cc.args.Parse()
			if err != nil {
				t.Errorf("not expected error: %q", err)
			}

			cc.data = tt.data

			result := cc.Cut()

			if len(result) == 0 && len(tt.expectedResult) == 0 {
				return
			}

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("got %v, expected %v", result, tt.expectedResult)
			}
		})
	}
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
