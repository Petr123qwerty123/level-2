package main

import (
	"flag"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"wb-level-2/develop/dev05/utils"
)

// Helper function to reset the command-line args.
func resetArgs(args []string) {
	os.Args = []string{"testArgs"}

	for _, arg := range args {
		os.Args = append(os.Args, arg)
	}
}

func TestGrepFlags_Parse(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		flags GrepFlags
	}{
		{
			name: "No flags",
			args: []string{},
			flags: GrepFlags{
				after:   -1,
				before:  -1,
				context: -1,
			},
		},
		{
			name: "After flag",
			args: []string{"-A", "5"},
			flags: GrepFlags{
				after:   5,
				before:  -1,
				context: -1,
			},
		},
		{
			name: "Before flag",
			args: []string{"-B", "3"},
			flags: GrepFlags{
				after:   -1,
				before:  3,
				context: -1,
			},
		},
		{
			name: "Context flag",
			args: []string{"-C", "2"},
			flags: GrepFlags{
				after:   -1,
				before:  -1,
				context: 2,
			},
		},
		{
			name: "Count flag",
			args: []string{"-c"},
			flags: GrepFlags{
				after:   -1,
				before:  -1,
				context: -1,
				count:   true,
			},
		},
		{
			name: "IgnoreCase flag",
			args: []string{"-i"},
			flags: GrepFlags{
				after:      -1,
				before:     -1,
				context:    -1,
				ignoreCase: true,
			},
		},
		{
			name: "Invert flag",
			args: []string{"-v"},
			flags: GrepFlags{
				after:   -1,
				before:  -1,
				context: -1,
				invert:  true,
			},
		},
		{
			name: "Fixed flag",
			args: []string{"-F"},
			flags: GrepFlags{
				after:   -1,
				before:  -1,
				context: -1,
				fixed:   true,
			},
		},
		{
			name: "LineNum flag",
			args: []string{"-n"},
			flags: GrepFlags{
				after:   -1,
				before:  -1,
				context: -1,
				lineNum: true,
			},
		},
		{
			name: "Invalid flag",
			args: []string{"-x"},
			flags: GrepFlags{
				after:   -1,
				before:  -1,
				context: -1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			gf := GrepFlags{}

			resetArgs(tt.args)
			gf.Parse()

			if gf != tt.flags {
				t.Errorf("GrepFlags.Parse() got = %v, want %v", gf, tt.flags)
			}
		})
	}
}

func TestGrepArgs_Parse(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		grepArgs GrepArgs
		hasErr   bool
	}{
		{
			name:     "No args",
			args:     []string{},
			grepArgs: GrepArgs{},
			hasErr:   true,
		},
		{
			name: "Pattern only",
			args: []string{"test"},
			grepArgs: GrepArgs{
				inputFiles: []*os.File{os.Stdin},
				pattern:    regexp.MustCompile("test"),
			},
			hasErr: false,
		},
		{
			name:     "Invalid pattern",
			args:     []string{"[invalid pattern"},
			grepArgs: GrepArgs{},
			hasErr:   true,
		},
		{
			name: "Invalid filename",
			args: []string{"test", "kjhg??fd"},
			grepArgs: GrepArgs{
				pattern: regexp.MustCompile("test"),
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			ga := GrepArgs{}

			resetArgs(tt.args)

			flag.Parse()
			err := ga.Parse()

			if (err != nil) != tt.hasErr {
				t.Errorf("GrepArgs.Parse() error = %v, wantErr %v", err, tt.hasErr)
				return
			}

			if !reflect.DeepEqual(ga.pattern, tt.grepArgs.pattern) {
				t.Errorf("GrepArgs.Parse() pattern = %v, want %v", ga.pattern, tt.grepArgs.pattern)
			}
			if len(ga.inputFiles) != len(tt.grepArgs.inputFiles) {
				t.Errorf("GrepArgs.Parse() inputFiles length = %v, want %v", len(ga.inputFiles), len(tt.grepArgs.inputFiles))
			}
		})
	}
}

func TestGrepClient_Grep(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		data           []string
		expectedResult []string
	}{
		{
			name: "Search with no flags",
			args: []string{"er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: []string{
				"raspberry",
				"watermelon",
				"blueberry",
			},
		},
		{
			name: "Search with after flag",
			args: []string{"-A", "3", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: []string{
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"blueberry",
			},
		},
		{
			name: "Search with before flag",
			args: []string{"-B", "3", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: []string{
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
		},
		{
			name: "Search with context flag",
			args: []string{"-C", "2", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: []string{
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"cheRry",
				"pomelo",
				"blueberry",
			},
		},
		{
			name: "Search with after and context flag",
			args: []string{"-A", "3", "-C", "2", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: []string{
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"cheRry",
				"pomelo",
				"blueberry",
			},
		},
		{
			name: "Search with before and context flag",
			args: []string{"-B", "3", "-C", "2", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: []string{
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"cheRry",
				"pomelo",
				"blueberry",
			},
		},
		{
			name: "Search with after, before with context flag",
			args: []string{"-A", "3", "-B", "4", "-C", "2", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: []string{
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"cheRry",
				"pomelo",
				"blueberry",
			},
		},
		{
			name: "Search with ignore case flag",
			args: []string{"-i", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: []string{
				"strawbErry",
				"raspberry",
				"watermelon",
				"cheRry",
				"blueberry",
			},
		},
		{
			name: "Search with invert flag",
			args: []string{"-v", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"grapes",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
			},
		},
		{
			name: "Search with fixed flag",
			args: []string{"-F", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: []string{},
		},
		{
			name: "Search with line num flag",
			args: []string{"-n", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: []string{
				"7. raspberry",
				"9. watermelon",
				"15. blueberry",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			resetArgs(tt.args)

			gc := GrepClient{}

			gc.flags.Parse()
			err := gc.args.Parse()

			gc.data = tt.data

			result, err := gc.Grep()

			if err != nil {
				t.Errorf("not expected error: %q", err)
			}

			if len(result) == 0 && len(tt.expectedResult) == 0 {
				return
			}

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("got %v, expected %v", result, tt.expectedResult)
			}
		})
	}
}

func TestGrepClient_Count(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		data           []string
		expectedResult int
	}{
		{
			name: "Count searched data with no flags",
			args: []string{"er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: 3,
		},
		{
			name: "Count searched data with after flag",
			args: []string{"-A", "3", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: 7,
		},
		{
			name: "Count searched data with before flag",
			args: []string{"-B", "3", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: 10,
		},
		{
			name: "Count searched data with context flag",
			args: []string{"-C", "2", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: 10,
		},
		{
			name: "Count searched data with after and context flag",
			args: []string{"-A", "3", "-C", "2", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: 10,
		},
		{
			name: "Count searched data with before and context flag",
			args: []string{"-B", "3", "-C", "2", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: 10,
		},
		{
			name: "Count searched data with after, before with context flag",
			args: []string{"-A", "3", "-B", "4", "-C", "2", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: 10,
		},
		{
			name: "Count searched data with ignore case flag",
			args: []string{"-i", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: 5,
		},
		{
			name: "Count searched data with invert flag",
			args: []string{"-v", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: 12,
		},
		{
			name: "Count searched data with fixed flag",
			args: []string{"-F", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: 0,
		},
		{
			name: "Count searched data with line num flag",
			args: []string{"-n", "er"},
			data: []string{
				"apple",
				"kiwi",
				"starfruit",
				"strawbErry",
				"pineapple",
				"mango",
				"raspberry",
				"grapes",
				"watermelon",
				"grapefruit",
				"banana",
				"pear",
				"cheRry",
				"pomelo",
				"blueberry",
			},
			expectedResult: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			resetArgs(tt.args)

			gc := GrepClient{}

			gc.flags.Parse()
			err := gc.args.Parse()

			gc.data = tt.data

			result, err := gc.Count()

			if err != nil {
				t.Errorf("not expected error: %q", err)
			}

			if result != tt.expectedResult {
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
