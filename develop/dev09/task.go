package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"wb-level-2/develop/dev09/utils"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var errNoArgs = errors.New("invalid number of arguments")

// WgetArgs структура, определяющая неименованные аргументы запуска утилиты Wget
type WgetArgs struct {
	urls []*url.URL
}

// Parse метод для распарсивания и сохранения значений неименованных аргументов запуска утилиты Wget в поля структуры
// WgetArgs
func (wa *WgetArgs) Parse() error {
	// так как в этой утилите по условию задачи непредусмотрено использование опций, парсим их в этом методе,
	// для получения неименованных аргументов
	flag.Parse()

	args := flag.Args()
	nArg := flag.NArg()

	// утилита имеет один обязательный неименованный аргумент - urls, если агрументов 0 - выводим ошибку, в любых других
	// случаях распарсиваем каждый переданный аргумент как url и складываем в слайс []*url.URL urls, если при парсе
	// url'а возникает ошибка - передаем её вверх по цепочке
	switch nArg {
	case 0:
		return errNoArgs
	default:
		for _, arg := range args {
			fileUrl, err := url.Parse(arg)
			if err != nil {
				wa.urls = []*url.URL{}
				return err
			}

			wa.urls = append(wa.urls, fileUrl)
		}
	}

	return nil
}

// WgetClient структура для управления утилитой Wget
type WgetClient struct {
	args WgetArgs
}

// NewWgetClient конструктор для создания объекта структуры WgetClient
func NewWgetClient() (*WgetClient, error) {
	// создание пустого объекта структуры WgetClient
	wc := &WgetClient{}

	// парс аргументов запуска утилиты
	err := wc.args.Parse()
	if err != nil {
		return nil, err
	}

	return wc, nil
}

// Wget метод для функционирования утилиты, который скачивает контент, находящийся по ссылкам, переданным в качестве
// неименованных аргументов утилиты. Чтобы не прирывать скачивание контента, находящимся по ссылкам, следующим за той,
// на которой была вызвана ошибка, мы не возвращаем ошибки сразу, а складываем их в слайс errs []error. Возвращает
// ошибки при скачивании.
func (wc *WgetClient) Wget() []error {
	var errs []error

	for _, fileUrl := range wc.args.urls {
		rawUrl := fileUrl.String()
		// получаем название файла для сохранения по ссылке
		fileName := utils.GetFileNameByUrl(fileUrl)

		// скачиваем контент, если это возможно, в ином случае - складываем ошибку в errs
		err := utils.DownloadContent(rawUrl, fileName)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// возвращаем отчет об ошибках
	return errs
}

// Start метод запуска утилиты
func (wc *WgetClient) Start() error {
	outputData := wc.Wget()

	err := utils.WriteData(os.Stderr, outputData...)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var wgetClient *WgetClient
	var err error

	// создание объекта структуры WgetClient, в случае ошибки - её вывод и выход из программы с кодом ошибки 1
	wgetClient, err = NewWgetClient()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}

	// запуск утилиты, в случае ошибки - её вывод и выход из программы с кодом ошибки 1
	err = wgetClient.Start()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}
}
