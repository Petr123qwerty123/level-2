package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"wb-level-2/develop/dev06/utils"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var errParse = errors.New("parse error")

type IntSlice []int

func (is *IntSlice) MarshalText() ([]byte, error) {
	return json.Marshal(*is)
}

func (is *IntSlice) UnmarshalText(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	str := string(b)

	parts := strings.Split(str, ",")

	for _, part := range parts {
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")

			if len(rangeParts) != 2 {
				return errParse
			}

			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return errParse
			}

			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return errParse
			}

			for i := start; i <= end; i++ {
				*is = append(*is, i)
			}
		} else {
			num, err := strconv.Atoi(part)
			if err != nil {
				return errParse
			}

			*is = append(*is, num)
		}
	}

	sort.Ints(*is)

	return nil
}

type CutFlags struct {
	fields    IntSlice
	delimiter string
	separated bool
}

func (cf *CutFlags) Parse() {
	flag.TextVar(&cf.fields, "f", &cf.fields, "Specify fields for cutting")
	flag.StringVar(&cf.delimiter, "d", "\t", "Specify another delimiter")
	flag.BoolVar(&cf.separated, "s", false, "Output lines only with delimiters")

	flag.Parse()
}

type CutArgs struct {
	inputFiles []*os.File
}

func (ca *CutArgs) Parse() error {
	args := flag.Args()
	nArg := flag.NArg()

	switch nArg {
	case 0:
		ca.inputFiles = append(ca.inputFiles, os.Stdin)
	default:
		for _, arg := range args {
			inputFile, err := os.Open(arg)
			if err != nil {
				return err
			}

			ca.inputFiles = append(ca.inputFiles, inputFile)
		}
	}

	return nil
}

type CutClient struct {
	flags CutFlags
	args  CutArgs
	data  []string
}

func NewCutClient() (*CutClient, error) {
	cc := &CutClient{}

	cc.flags.Parse()
	err := cc.args.Parse()
	if err != nil {
		return nil, err
	}

	for _, inputFile := range cc.args.inputFiles {
		partData, err := utils.ReadData(inputFile)
		if err != nil {
			_ = inputFile.Close()
			return nil, err
		}

		err = inputFile.Close()
		if err != nil {
			return nil, err
		}

		cc.data = append(cc.data, partData...)
	}

	return cc, nil
}

func (cc *CutClient) Cut() ([]string, error) {
	var result []string

	for _, line := range cc.data {
		fields := strings.Split(line, cc.flags.delimiter)

		if !(cc.flags.separated && len(fields) == 1) {
			var builder strings.Builder

			for i, fieldNumber := range cc.flags.fields {
				i := i

				if fieldNumber >= 1 && fieldNumber <= len(fields) {
					builder.WriteString(fields[fieldNumber-1])

					if i != len(cc.flags.fields)-1 {
						builder.WriteString(cc.flags.delimiter)
					}
				}
			}

			if len(cc.flags.fields) > 0 {
				result = append(result, builder.String())
			} else {
				result = append(result, line)
			}
		}
	}

	return result, nil
}

func (cc *CutClient) Start() error {
	outputData, err := cc.Cut()

	err = utils.WriteData(os.Stdout, outputData...)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var cutClient *CutClient
	var err error

	cutClient, err = NewCutClient()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}

	err = cutClient.Start()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}
}
