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

type WgetArgs struct {
	urls []*url.URL
}

func (wa *WgetArgs) Parse() error {
	flag.Parse()

	args := flag.Args()
	nArg := flag.NArg()

	switch nArg {
	case 0:
		return errors.New("invalid number of arguments")
	default:
		for _, arg := range args {
			fileUrl, err := url.Parse(arg)
			if err != nil {
				return err
			}

			wa.urls = append(wa.urls, fileUrl)
		}
	}

	return nil
}

type WgetClient struct {
	args WgetArgs
}

func NewWgetClient() (*WgetClient, error) {
	wc := &WgetClient{}

	err := wc.args.Parse()
	if err != nil {
		return nil, err
	}

	return wc, nil
}

func (wc *WgetClient) Wget() []error {
	var errs []error

	for _, fileUrl := range wc.args.urls {
		rawUrl := fileUrl.String()
		fileName := utils.GetFileNameByUrl(fileUrl)

		err := utils.DownloadContent(rawUrl, fileName)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

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

	wgetClient, err = NewWgetClient()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}

	err = wgetClient.Start()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}
}
