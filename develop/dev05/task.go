package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"wb-level-2/develop/dev05/utils"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения +
-B - "before" печатать +N строк до совпадения +
-C - "context" (A+B) печатать ±N строк вокруг совпадения +
-c - "count" (количество строк) +
-i - "ignore-case" (игнорировать регистр) +
-v - "invert" (вместо совпадения, исключать) +
-F - "fixed", точное совпадение со строкой, не паттерн +
-n - "line num", печатать номер строки +

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type GrepFlags struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
}

func (gf *GrepFlags) Parse() {
	flag.IntVar(&gf.after, "A", -1, "Output +N lines after match")
	flag.IntVar(&gf.before, "B", -1, "Output +N lines before match")
	flag.IntVar(&gf.context, "C", -1, "Output ±N lines around the match")
	flag.BoolVar(&gf.count, "c", false, "Output the number of lines found")
	flag.BoolVar(&gf.ignoreCase, "i", false, "Ignore case characters")
	flag.BoolVar(&gf.invert, "v", false, "Output lines in which the search pattern is not found")
	flag.BoolVar(&gf.fixed, "F", false, "Use a regular string instead of a regular expression")
	flag.BoolVar(&gf.lineNum, "n", false, "Output line number with the found lines")

	flag.Parse()
}

type GrepArgs struct {
	pattern    *regexp.Regexp
	inputFiles []*os.File
}

func (ga *GrepArgs) Parse() error {
	var err error

	args := flag.Args()
	nArg := flag.NArg()

	switch nArg {
	case 1:
		pattern := flag.Arg(0)

		ga.pattern, err = regexp.Compile(pattern)
		if err != nil {
			return err
		}

		ga.inputFiles = append(ga.inputFiles, os.Stdin)
	case 2:
		pattern := flag.Arg(0)
		filenames := args[1:]

		ga.pattern, err = regexp.Compile(pattern)
		if err != nil {
			return err
		}

		for _, filename := range filenames {
			inputFile, err := os.Open(filename)
			if err != nil {
				return err
			}

			ga.inputFiles = append(ga.inputFiles, inputFile)
		}
	default:
		return errors.New("invalid number of arguments")
	}

	return nil
}

type GrepClient struct {
	flags GrepFlags
	args  GrepArgs
	data  []string
}

func NewGrepClient() (*GrepClient, error) {
	gc := &GrepClient{}

	gc.flags.Parse()
	err := gc.args.Parse()
	if err != nil {
		return nil, err
	}

	for _, inputFile := range gc.args.inputFiles {
		partData, err := utils.ReadData(inputFile)
		if err != nil {
			_ = inputFile.Close()
			return nil, err
		}

		err = inputFile.Close()
		if err != nil {
			return nil, err
		}

		gc.data = append(gc.data, partData...)
	}

	if gc.flags.ignoreCase {
		gc.args.pattern, err = regexp.Compile("(?i)" + gc.args.pattern.String())
		if err != nil {
			return nil, err
		}
	}

	if gc.flags.fixed {
		gc.args.pattern, err = regexp.Compile("^" + gc.args.pattern.String() + "$")
		if err != nil {
			return nil, err
		}
	}

	return gc, nil
}

func (gc *GrepClient) Grep() ([]string, error) {
	var after, before int
	var result []string

	lenData := len(gc.data)
	usingIndexes := make([]bool, lenData)

	if gc.flags.context > 0 {
		after, before = gc.flags.context, gc.flags.context
	} else {
		if gc.flags.after > 0 {
			after = gc.flags.after
		}

		if gc.flags.before > 0 {
			before = gc.flags.before
		}
	}

	if gc.flags.invert {
		for i := 0; i < lenData; i++ {
			usingIndexes[i] = !usingIndexes[i]
		}
	}

	for i, str := range gc.data {
		if gc.args.pattern.MatchString(str) {
			for j := i - before; j < i; j++ {
				j := j

				if j >= 0 && !usingIndexes[j] {
					usingIndexes[j] = true
				}
			}

			if !usingIndexes[i] {
				usingIndexes[i] = true
			}

			for j := i + 1; j < i+after+1; j++ {
				j := j

				if j <= lenData-1 && !usingIndexes[j] {
					usingIndexes[j] = true
				}
			}
		}
	}

	for index, use := range usingIndexes {
		if use {
			if gc.flags.lineNum {
				result = append(result, strconv.Itoa(index+1)+". "+gc.data[index])
			} else {
				result = append(result, gc.data[index])
			}
		}
	}

	return result, nil
}

func (gc *GrepClient) Count() (int, error) {
	searchResult, err := gc.Grep()
	if err != nil {
		return -1, err
	}

	lenSearchResult := len(searchResult)

	return lenSearchResult, nil
}

func (gc *GrepClient) Start() error {
	if gc.flags.count {
		outputData, err := gc.Count()
		err = utils.WriteData(os.Stdout, outputData)
		if err != nil {
			return err
		}
	} else {
		outputData, err := gc.Grep()
		err = utils.WriteData(os.Stdout, outputData...)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var grepClient *GrepClient
	var err error

	grepClient, err = NewGrepClient()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}

	err = grepClient.Start()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}
}
