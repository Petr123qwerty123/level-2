package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ReadData принимает на вход reader, возвращает слайс прочитанных строк
func ReadData(reader io.Reader) ([]string, error) {
	var lines []string

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// WriteData принимает на вход writer и слайс данных любого типа data для записи, записывает построчно data в writer
func WriteData[T any](writer io.Writer, data ...T) error {
	for _, line := range data {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			return err
		}
	}

	return nil
}

// TrimLeadingSpaces убирает у s слева любые пробельные символы, возвращает мутированную строку
func TrimLeadingSpaces(s string) string {
	return strings.TrimLeftFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

// ParseMonth принимает на вхож строку, возвращает индекс месяца, если, убрав первые пробельные символы, первые три
// символа окажутся началом названия месяца на английском, иначе -1 и ошибку
func ParseMonth(word string) (int, error) {
	var errParse = errors.New("month cannot be parsed")
	// убираем первые пробельные символы
	word = TrimLeadingSpaces(word)

	// если символов в оставшейся строке меньше 3, то этого количества символов недостаточно для однозначного
	// определения месяца
	if len(word) < 3 {
		return -1, errParse
	}

	// первые три символа в верхнем регистре
	firstThreeLetters := strings.ToTitle(word[:3])

	switch firstThreeLetters {
	case "JAN":
		return 0, nil
	case "FEB":
		return 1, nil
	case "MAR":
		return 2, nil
	case "APR":
		return 3, nil
	case "MAY":
		return 4, nil
	case "JUN":
		return 5, nil
	case "JUL":
		return 6, nil
	case "AUG":
		return 7, nil
	case "SEP":
		return 8, nil
	case "OCT":
		return 9, nil
	case "NOV":
		return 10, nil
	case "DEC":
		return 11, nil
	default:
		return -1, errParse
	}
}

// ParseNumericValue принимает на вход строку и флаг об использовании числового суффикса, возвращает распаршенное число
// типа float64, если убрав первые пробельные символы, начало строки будет соответсвовать регулярному выражению
// ^\d+(\.\d+)?[KkMmGgTt]?, иначе -1 и ошибку
func ParseNumericValue(word string, useSuffix bool) (float64, error) {
	var errParse = errors.New("number cannot be parsed")
	// убираем первые пробельные символы
	word = TrimLeadingSpaces(word)

	// регулярное выражение соответсвия числу, в том числе с числовым суффиксом
	pattern := `^\d+(\.\d+)?[KkMmGgTt]?`
	// компиляция регулярного выражения
	regex, err := regexp.Compile(pattern)

	if err != nil {
		return -1, err
	}

	// поиск числа в строке по регулярному выражению
	word = regex.FindString(word)

	// если число не было найдено возвращает -1 и ошибку
	if len(word) == 0 {
		return -1, errParse
	}

	// если нужно было учитывать суффикс, а суффикса на самом деле нет добавляем заглушку _ в конец строки и парсим,
	// иначе если суффикс не нужно учитывать, а он есть - убираем и парсим
	if useSuffix {
		if unicode.IsDigit(rune(word[len(word)-1])) {
			word = word + "_"
		}

		suffix := strings.ToTitle(string(word[len(word)-1]))

		num, err := strconv.ParseFloat(word[:len(word)-1], 64)
		if err != nil {
			return 0, err
		}

		switch suffix {
		case "K":
			num *= 1000
		case "M":
			num *= 1000000
		case "G":
			num *= 1000000000
		case "T":
			num *= 1000000000000
		}

		return num, nil
	} else {
		if unicode.IsDigit(rune(word[len(word)-1])) {
			return strconv.ParseFloat(word, 64)
		}

		return strconv.ParseFloat(word[:len(word)-1], 64)
	}
}

// RemoveDuplicates удаляет дубликаты в слайсе строк с помощью мапы
func RemoveDuplicates(lines []string) []string {
	var result []string

	uniqueLines := make(map[string]bool)

	for _, line := range lines {
		if !uniqueLines[line] {
			uniqueLines[line] = true
			result = append(result, line)
		}
	}
	return result
}
