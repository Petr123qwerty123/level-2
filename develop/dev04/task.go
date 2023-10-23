package main

import (
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func ToLowerUnique(words []string) []string {
	uniqueWords := make(map[string]bool)

	result := make([]string, 0)

	for _, word := range words {
		lowerWord := strings.ToLower(word)

		if !uniqueWords[lowerWord] {
			uniqueWords[lowerWord] = true
			result = append(result, lowerWord)
		}
	}

	return result
}

func SortWord(word string) string {
	wordArr := strings.Split(word, "")

	sort.Strings(wordArr)

	return strings.Join(wordArr, "")
}

func FindAnagramSets(words []string) map[string][]string {
	anagramSets := make(map[string][]string)

	keys := make(map[string]string)

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		sortedLowerWord := SortWord(lowerWord)

		_, exists := keys[sortedLowerWord]
		if !exists {
			keys[sortedLowerWord] = word
		}
	}

	words = ToLowerUnique(words)

	for _, word := range words {
		sortedWord := SortWord(word)

		anagramSets[sortedWord] = append(anagramSets[sortedWord], word)
	}

	for key, value := range keys {
		anagramSets[value] = anagramSets[key]
		delete(anagramSets, key)
	}

	for sortedWord, wordSet := range anagramSets {
		if len(wordSet) <= 1 {
			delete(anagramSets, sortedWord)
		} else {
			sort.Strings(wordSet)
		}
	}

	return anagramSets
}
