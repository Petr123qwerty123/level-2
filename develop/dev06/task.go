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

// IntSlice тип для задания слайса по списку точек и отрезков. Точка - целое число, отрезок - множество целых чисел,
// находящееся между крайними точками отрезка, включая сами крайние точки. Перечисление элементов (точек и отрезков)
// осуществляется с помощью запятой
type IntSlice []int

// MarshalText метод для сериализации слайса целых чисел
func (is *IntSlice) MarshalText() ([]byte, error) {
	return json.Marshal(*is)
}

// UnmarshalText метод для десериализации в слайс целых чисел
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

// CutFlags структура, определяющюая опции утилиты Cut
type CutFlags struct {
	fields    IntSlice
	delimiter string
	separated bool
}

// Parse метод для распарсивания и сохранения значений флагов опций в поля структуры CutFlags
func (cf *CutFlags) Parse() {
	flag.TextVar(&cf.fields, "f", &cf.fields, "Specify fields for cutting")
	flag.StringVar(&cf.delimiter, "d", "\t", "Specify another delimiter")
	flag.BoolVar(&cf.separated, "s", false, "Output lines only with delimiters")

	flag.Parse()
}

// CutArgs структура, определяющая неименованные аргументы запуска утилиты Cut
type CutArgs struct {
	inputFiles []*os.File
}

// Parse метод для распарсивания и сохранения значений неименованных аргументов запуска утилиты Cut в поля структуры
// CutArgs
func (ca *CutArgs) Parse() error {
	args := flag.Args()
	nArg := flag.NArg()

	// если количество неименованных аргументов (источник данных для сортировки) - 0, то добавляем в источник os.Stdin,
	// иначе воспринимаем введенные агрументы как пути до файлов, с которых нужно будет считывать данные, открываем их
	// и кладем объекты *os.File в inputFiles
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

// CutClient структура для управления утилитой Cut
type CutClient struct {
	flags CutFlags
	args  CutArgs
	data  []string
}

// NewCutClient конструктор для создания объекта структуры CutClient
func NewCutClient() (*CutClient, error) {
	// создание пустого объекта структуры CutClient
	cc := &CutClient{}

	// парс флагов и аргументов запуска утилиты
	cc.flags.Parse()
	err := cc.args.Parse()
	if err != nil {
		return nil, err
	}

	// чтение и сохранение данных для сортировки в поле data структуры CutClient, закрытие reader'ов
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

// Cut метод для функционирования утилиты, который возвращает вырезанные данные, исходя из установленных опций
// при запуске утилиты, в случае отсутствия ошибок
func (cc *CutClient) Cut() []string {
	var result []string

	// проходимся циклом по строкам данных для обработки
	for _, line := range cc.data {
		// делим строку на поля разделителем
		fields := strings.Split(line, cc.flags.delimiter)

		// в случае если стоит флаг -s, указывающий на использование только строк с разделителем, а при сплите строки
		// разделителем мы получили слайс длиной один - это значит, что разделителя в строке нет, значит нам ничего
		// обрабатывать и выводить эту строку не нужно, значит мы работаем при обратном условии
		if !(cc.flags.separated && len(fields) == 1) {
			var builder strings.Builder

			// проходимся циклом по слайсу номеров колонок, которые вырезаются
			for i, fieldNumber := range cc.flags.fields {
				// переопределяем индекс итерации внутрии цикла
				i := i

				// если номер строки не выходит за границы слайса - обрабатываем
				if fieldNumber >= 1 && fieldNumber <= len(fields) {
					// если поле в строке не последнее, то записываем поле + разделитель, в ином случае просто поле
					builder.WriteString(fields[fieldNumber-1])

					if i != len(cc.flags.fields)-1 {
						builder.WriteString(cc.flags.delimiter)
					}
				}
			}

			// если был передан флаг -f, то склеиваем эти поля, если нет - в результат записываем всю строку
			if len(cc.flags.fields) > 0 {
				result = append(result, builder.String())
			} else {
				result = append(result, line)
			}
		}
	}

	return result
}

// Start метод запуска утилиты
func (cc *CutClient) Start() error {
	outputData := cc.Cut()

	err := utils.WriteData(os.Stdout, outputData...)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var cutClient *CutClient
	var err error

	// создание объекта структуры CutClient, в случае ошибки - её вывод и выход из программы с кодом ошибки 1
	cutClient, err = NewCutClient()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}

	// запуск утилиты, в случае ошибки - её вывод и выход из программы с кодом ошибки 1
	err = cutClient.Start()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}
}
