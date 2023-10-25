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

func WriteData[T any](writer io.Writer, data ...T) error {
	for _, line := range data {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			return err
		}
	}

	return nil
}

func TrimLeadingSpaces(s string) string {
	return strings.TrimLeftFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

func ParseMonth(word string) (int, error) {
	word = TrimLeadingSpaces(word)

	if len(word) < 3 {
		return -1, errors.New("month cannot be parsed")
	}

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
		return -1, errors.New("month cannot be parsed")
	}
}

func ParseNumericValue(word string, suffix bool) (float64, error) {
	word = TrimLeadingSpaces(word)

	pattern := `^\d+(\.\d+)?[KkMmGgTt]?`

	regex, err := regexp.Compile(pattern)

	if err != nil {
		return -1, err
	}

	word = regex.FindString(word)

	if len(word) == 0 {
		return -1, errors.New("number cannot be parsed")
	}

	if suffix {
		if unicode.IsDigit(rune(word[len(word)-1])) {
			word = word + "_"
		}

		num, err := strconv.ParseFloat(word[:len(word)-1], 64)
		if err != nil {
			return 0, err
		}

		switch word[len(word)-1] {
		case 'K' | 'k':
			num *= 1000
		case 'M' | 'm':
			num *= 1000000
		case 'G' | 'g':
			num *= 1000000000
		case 'T' | 't':
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
