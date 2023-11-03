package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func WriteData[T any](writer io.Writer, data ...T) error {
	for _, line := range data {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetFileNameByUrl принимает на вход url, возвращает название файла для сохранения контента, находящегося по ссылке
func GetFileNameByUrl(fileUrl *url.URL) string {
	path := fileUrl.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	if fileName == "" {
		fileName = "index.html"
	}

	return fileName
}

// DownloadContent принимает на вход ссылку в виде string и название файла, под которым сохранить скачанный контент,
// находящийся по ссылке
func DownloadContent(rawUrl, fileName string) error {
	// делаем get-запрос по переданной ссылке
	resp, err := http.Get(rawUrl)
	// если resp != nil, значит ошибки нет, значит мы сможем закрыть resp.Body, поэтому в defer не отлавливаем ошибку,
	// если resp == nil, значит есть ошибка, значит мы её вернем ниже
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	// если статус ответа не 200 - возвращаем ошибку
	if resp.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}

	// в ином случае создаем файл c названием, переданным в функцию
	file, err := os.Create(fileName)
	// по аналогии, если file != nil, значит ошибки нет, значит мы сможем закрыть file, поэтому в defer не отлавливаем ошибку,
	// если file == nil, значит есть ошибка, значит мы её вернем ниже
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return err
	}

	// переносим весь контент из resp.Body в file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
