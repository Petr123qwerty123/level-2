package unpacker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// PackedChar структура запакованного символа, где ch - символ, nr (number of repetitions) - количество повторений
// символа от 0 до 9
type PackedChar struct {
	ch rune
	nr int // nr [0, 9]
}

// NewPackedChar конструктор PackedChar
// возвращает объект структуры PackedChar в случае, если количество повторений - цифра (число от 0 до 9),
// в ином случае - ошибку
func NewPackedChar(ch rune, nr int) (*PackedChar, error) {
	if nr < 0 || nr > 9 {
		return nil, fmt.Errorf("invalid repetition number: %d", nr)
	}

	return &PackedChar{ch: ch, nr: nr}, nil
}

// Unpack распаковывает символ (возвращает строку, где ch повторяется nr раз)
func (pc PackedChar) Unpack() string {
	return strings.Repeat(string(pc.ch), pc.nr)
}

// PackedString - тип запакованной строки (слайс запакованных символов)
type PackedString []PackedChar

// NewPackedString конструктор PackedString
// на вход принимает строку s, которая будет проверена на правильность формата. В случае правильного формата вернется
// объект PackedString, в ином - ошибка
func NewPackedString(s string) (*PackedString, error) {
	var packedString PackedString

	runes := []rune(s)

	var ch rune
	var nr int
	toDelete := 1

	for len(runes) != 0 {
		if runes[0] >= 48 && runes[0] <= 57 {
			return nil, errors.New("invalid string")
		} else if runes[0] == 92 {
			if len(runes) >= 2 && (runes[1] == 92 || (runes[1] >= 48 && runes[1] <= 57)) {
				ch = runes[1]
				nr = 1
				toDelete = 2

				if len(runes) >= 3 && runes[2] >= 48 && runes[2] <= 57 {
					nr, _ = strconv.Atoi(string(runes[2]))
					toDelete = 3
				}
			} else {
				return nil, errors.New("invalid string")
			}
		} else {
			ch = runes[0]
			nr = 1
			toDelete = 1
			if len(runes) >= 2 && runes[1] >= 48 && runes[1] <= 57 {
				nr, _ = strconv.Atoi(string(runes[1]))
				toDelete = 2
			}
		}

		packedChar, _ := NewPackedChar(ch, nr)
		packedString = append(packedString, *packedChar)
		runes = runes[toDelete:]
	}

	return &packedString, nil
}

// Unpack распаковывает строку (возвращает строку, где каждый запакованный символ будет распакован)
func (ps PackedString) Unpack() string {
	var builder strings.Builder

	for _, pch := range ps {
		builder.WriteString(pch.Unpack())
	}

	return builder.String()
}
