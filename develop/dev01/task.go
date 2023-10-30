package main

import (
	"github.com/beevik/ntp"
	"time"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

// PrintCurrentTime возвращает время с ntp сервера
func PrintCurrentTime() (time.Time, error) {
	// запрос к ntp серверу
	response, err := ntp.Query("0.ru.pool.ntp.org")
	// достаем с response время и добавляем минимальную ошибку между временем клиента (нас) и сервером
	ntpTime := response.Time.Add(response.MinError)

	return ntpTime, err
}
