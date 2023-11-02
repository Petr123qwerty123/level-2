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

var errNoArgs = errors.New("invalid number of arguments")

// GrepFlags структура, определяющюая опции утилиты Grep
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

// Parse метод для распарсивания и сохранения значений флагов опций в поля структуры GrepFlags
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

// GrepArgs структура, определяющая неименованные аргументы запуска утилиты Grep
type GrepArgs struct {
	pattern    *regexp.Regexp
	inputFiles []*os.File
}

// Parse метод для распарсивания и сохранения значений неименованных аргументов запуска утилиты Grep в поля структуры
// GrepArgs
func (ga *GrepArgs) Parse() error {
	var err error

	args := flag.Args()
	nArg := flag.NArg()

	// утилита имеет один обязательный неименованный аргумент - pattern, если агрументов 0 - выводим ошибку, если
	// передан только один аргумент, воспринимаем его как pattern (регулярное выражение), в качестве источника поиска
	// данных берем os.Stdin, в случае передачи более 1 аргумента воспринимаме первый как pattern, остальные - как
	// пути до файлов (источников для поиска), открываем их и кладём *os.File в inputFiles
	switch nArg {
	case 0:
		return errNoArgs
	case 1:
		pattern := flag.Arg(0)

		ga.pattern, err = regexp.Compile(pattern)
		if err != nil {
			return err
		}

		ga.inputFiles = append(ga.inputFiles, os.Stdin)
	default:
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
	}

	return nil
}

// GrepClient структура для управления утилитой Grep
type GrepClient struct {
	flags GrepFlags
	args  GrepArgs
	data  []string
}

// NewGrepClient конструктор для создания объекта структуры GrepClient
func NewGrepClient() (*GrepClient, error) {
	// создание пустого объекта структуры GrepClient
	gc := &GrepClient{}

	// парс флагов и аргументов запуска утилиты
	gc.flags.Parse()
	err := gc.args.Parse()
	if err != nil {
		return nil, err
	}

	// чтение и сохранение данных для поиска в поле data структуры GrepClient, закрытие reader'ов
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

	return gc, nil
}

// Grep метод для функционирования утилиты, который возвращает найденные по паттерну данные, исходя из установленных
// опций при запуске утилиты
func (gc *GrepClient) Grep() ([]string, error) {
	var after, before int
	var result []string
	var err error

	// учитывание опции -i, редактирование с учетом этой опции регулярного выражения
	if gc.flags.ignoreCase {
		gc.args.pattern, err = regexp.Compile("(?i)" + gc.args.pattern.String())
		if err != nil {
			return nil, err
		}
	}

	// учитывание опции -F, редактирование с учетом этой опции регулярного выражения
	if gc.flags.fixed {
		gc.args.pattern, err = regexp.Compile("^" + gc.args.pattern.String() + "$")
		if err != nil {
			return nil, err
		}
	}

	// длина слайса с данными для поиска
	lenData := len(gc.data)
	// создаем слайс []bool длиной, равной слайсу с данными для поиска, чтобы понимать по индексу слайса, какую строку
	// из gc.data использовать, по умолчанию мы ничего не находим, то есть каждый элемент usingIndexes - false
	usingIndexes := make([]bool, lenData)

	// учитывание опций -A, -B, -C
	// флаг -C в приоритете, если есть значение у -C, то after=before=gc.flags.context,
	// если нет значение у флага -C, но есть у -A или -B, то after=gc.flags.after или before=gc.flags.before
	// соответсвенно
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

	// учитывание опции -v
	// так как ниже мы ищем данные и меняем значения в usingIndexes на противоположные (true->false, false->true),
	// то когда выставлен флаг -v мы сразу инвертируем каждый элемент usingIndexes (сразу все считается результатом),
	// соответсвенно когда мы действительно что-то найдем по паттерну, мы изменим значение на противоположное, таким
	// образом добьемся инвертации результата поиска (результат - все, кроме того, что найдено)
	if gc.flags.invert {
		for i := 0; i < lenData; i++ {
			usingIndexes[i] = !usingIndexes[i]
		}
	}

	// осуществление поиска нужных строк по паттерну с помощью итерации циклом по слайсу данных для поиска с учетом
	// добавления строк до и после в соответствии флагам -A, -B, -C
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

			if usingIndexes[i] && gc.flags.invert {
				usingIndexes[i] = false
			}

			for j := i + 1; j < i+after+1; j++ {
				j := j

				if j <= lenData-1 && !usingIndexes[j] {
					usingIndexes[j] = true
				}
			}
		}
	}

	// с помощью слайса usingIndexes мы в result добавляем только те элементы gc.data, значение по индексу которых в
	// слайсе usingIndexes - true
	for index, use := range usingIndexes {
		if use {
			// учитывание опции -n
			// если выставлен этот флаг, то к строке, добавляемой в result, прибавляется в начале её номер (индекс + 1)
			// и ". " как сепаратор между номером и самой строкой
			if gc.flags.lineNum {
				result = append(result, strconv.Itoa(index+1)+". "+gc.data[index])
			} else {
				result = append(result, gc.data[index])
			}
		}
	}

	return result, nil
}

// Count возвращает количество найденных строк по паттерну или ошибку
func (gc *GrepClient) Count() (int, error) {
	searchResult, err := gc.Grep()
	if err != nil {
		return -1, err
	}

	lenSearchResult := len(searchResult)

	return lenSearchResult, nil
}

func (gc *GrepClient) Start() error {
	// учитывание опции -c
	// в вывод идут не сами найденные строки, а количество таких строк
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

	// создание объекта структуры GrepClient, в случае ошибки - её вывод и выход из программы с кодом ошибки 1
	grepClient, err = NewGrepClient()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}

	// запуск утилиты, в случае - её вывод и выход из программы с кодом ошибки 1
	err = grepClient.Start()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}
}
