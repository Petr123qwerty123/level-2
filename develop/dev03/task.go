package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	s "sort"
	"strconv"
	"strings"
	"wb-level-2/develop/dev03/utils"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению +
-r — сортировать в обратном порядке +
-u — не выводить повторяющиеся строки +

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца +
-b — игнорировать хвостовые пробелы +
-c — проверять отсортированы ли данные +
-h — сортировать по числовому значению с учётом суффиксов +

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

	s.Ints(*is)

	return nil
}

type SortFlags struct {
	columns       IntSlice
	numeric       bool
	reverse       bool
	unique        bool
	month         bool
	ignoreSpaces  bool
	checkSorted   bool
	numericSuffix bool
}

func (sf *SortFlags) Parse() {
	flag.TextVar(&sf.columns, "k", &sf.columns, "Specify columns for sorting")
	flag.BoolVar(&sf.numeric, "n", false, "Sort by numeric value")
	flag.BoolVar(&sf.reverse, "r", false, "Sort in reverse order")
	flag.BoolVar(&sf.unique, "u", false, "Do not output repeated lines")
	flag.BoolVar(&sf.month, "M", false, "Sort by month name")
	flag.BoolVar(&sf.ignoreSpaces, "b", false, "Ignore trailing spaces")
	flag.BoolVar(&sf.checkSorted, "c", false, "Check if the data is sorted")
	flag.BoolVar(&sf.numericSuffix, "h", false, "Sort by numeric value with suffixes")

	flag.Parse()
}

type SortArgs struct {
	inputFiles []*os.File
}

func (sa *SortArgs) Parse() error {
	args := flag.Args()
	nArg := flag.NArg()

	switch nArg {
	case 0:
		sa.inputFiles = append(sa.inputFiles, os.Stdin)
	default:
		for _, arg := range args {
			inputFile, err := os.Open(arg)
			if err != nil {
				return err
			}

			sa.inputFiles = append(sa.inputFiles, inputFile)
		}
	}

	return nil
}

type SortClient struct {
	flags SortFlags
	args  SortArgs
	data  []string
}

func NewSortClient() (*SortClient, error) {
	sc := &SortClient{}

	sc.flags.Parse()
	err := sc.args.Parse()
	if err != nil {
		return nil, err
	}

	for _, inputFile := range sc.args.inputFiles {
		partData, err := utils.ReadData(inputFile)
		if err != nil {
			_ = inputFile.Close()
			return nil, err
		}

		err = inputFile.Close()
		if err != nil {
			return nil, err
		}

		sc.data = append(sc.data, partData...)
	}

	return sc, nil
}

func (sc *SortClient) Sort() ([]string, error) {
	result := make([]string, len(sc.data), cap(sc.data))

	copy(result, sc.data)

	slices.SortFunc(result, func(line1, line2 string) int {
		if len(sc.flags.columns) > 0 {
			fields1 := strings.Fields(line1)
			fields2 := strings.Fields(line2)

			var builder1 strings.Builder
			var builder2 strings.Builder

			for _, fieldNumber := range sc.flags.columns {
				if fieldNumber >= 1 && fieldNumber <= len(fields1) {
					builder1.WriteString(fields1[fieldNumber-1])
				}

				if fieldNumber >= 1 && fieldNumber <= len(fields2) {
					builder2.WriteString(fields2[fieldNumber-1])
				}
			}

			line1 = builder1.String()
			line2 = builder2.String()
		}

		if sc.flags.ignoreSpaces {
			line1 = utils.TrimLeadingSpaces(line1)
			line2 = utils.TrimLeadingSpaces(line2)
		}

		if sc.flags.numeric || sc.flags.numericSuffix {
			num1, err1 := utils.ParseNumericValue(line1, sc.flags.numericSuffix)
			num2, err2 := utils.ParseNumericValue(line2, sc.flags.numericSuffix)

			if err1 != nil || err2 != nil {
				return cmp.Compare(line1, line2)
			}

			return cmp.Compare(num1, num2)
		}

		if sc.flags.month {
			indMonth1, err1 := utils.ParseMonth(line1)
			indMonth2, err2 := utils.ParseMonth(line2)

			if err1 != nil || err2 != nil {
				return cmp.Compare(line1, line2)
			}

			return cmp.Compare(indMonth1, indMonth2)
		}

		return cmp.Compare(line1, line2)
	})

	if sc.flags.unique {
		result = utils.RemoveDuplicates(result)
	}

	if sc.flags.reverse {
		slices.Reverse(result)
	}

	return result, nil
}

func (sc *SortClient) IsSorted() (int, error) {
	sortedData, err := sc.Sort()
	if err != nil {
		return -1, err
	}

	high := min(len(sortedData), len(sc.data))

	for i := 0; i < high; i++ {
		if sc.data[i] != sortedData[i] {
			return i, nil
		}
	}

	return -1, nil
}

func (sc *SortClient) Start() error {
	if sc.flags.checkSorted {
		outputData, err := sc.IsSorted()
		err = utils.WriteData(os.Stdout, outputData)
		if err != nil {
			return err
		}
	} else {
		outputData, err := sc.Sort()
		err = utils.WriteData(os.Stdout, outputData...)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var sortClient *SortClient
	var err error

	sortClient, err = NewSortClient()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}

	err = sortClient.Start()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}
}
