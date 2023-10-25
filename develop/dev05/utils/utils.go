package utils

import (
	"bufio"
	"fmt"
	"io"
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
