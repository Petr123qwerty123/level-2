package main

import (
	"fmt"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“done after %v”, time.Since(start))
*/

/*Стоит сделать пометку о том, что задание сформулировано непонятно, поэтому функция реализована таким образом, что пока
открыт хотя бы один канал - or тоже открыт*/

// sig функция, реализация которой взята из условия к задаче
// принимает на вход продолжительность работы возвращаемого канала (время, через которое возвращаемый канал закроется)
func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})

	go func() {
		defer close(c)

		time.Sleep(after)
	}()

	return c
}

// or принимает на вход слайс каналов любого типа в режиме чтения, возвращает канал, объединяющий потоки чтения из
// каналов, переданных в функцию
func or(channels ...<-chan interface{}) <-chan interface{} {
	// создаем канал любого типа для объединения потоков чтения переданных в функцию каналов
	done := make(chan interface{})

	// определяем sync.WaitGroup для ожидания окончания работы всех воркеров, запущенных в цикле для передачи данных с
	// переданных в функцию каналов в общий
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	// запуск воркеров
	for _, channel := range channels {
		go func(channel <-chan interface{}, done chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()

			for value := range channel {
				done <- value
			}

			fmt.Println("channel closed")
		}(channel, done, &wg)
	}

	// запуск горутины, ожидающей, когда все воркеры остановятся, и закрывающий общий канал
	go func(done chan interface{}, wg *sync.WaitGroup) {
		wg.Wait()
		close(done)
	}(done, &wg)

	return done
}

func main() {
	start := time.Now()

	<-or(
		sig(10*time.Second),
		sig(9*time.Second),
		sig(11*time.Second),
	)

	fmt.Printf("done after %v", time.Since(start).Seconds())
}
