package main

import (
	"flag"
	"net/url"
	"os"
	"strings"
	"testing"
	"wb-level-2/develop/dev09/utils"
)

// Helper function which checks the existence of the downloaded file in the current directory
func checkExistFileInCurDir(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}

	return true
}

// Helper function to reset the command-line args.
func resetArgs(args []string) {
	os.Args = []string{"testArgs"}

	for _, arg := range args {
		os.Args = append(os.Args, arg)
	}
}

// Helper function to get slice of *url.URl by strings not taking into account possible errors
func urlMustParse(rawUrls ...string) []*url.URL {
	var result []*url.URL

	for _, rawUrl := range rawUrls {
		u, _ := url.Parse(rawUrl)
		result = append(result, u)
	}

	return result
}

func TestGetFileNameByUrl(t *testing.T) {
	tests := []struct {
		name           string
		rawUrl         string
		expectedResult string
	}{
		{
			name:           "Get filename by simple url",
			rawUrl:         "https://example.com/files/document.pdf",
			expectedResult: "document.pdf",
		},
		{
			name:           "Get filename by url without filename inside",
			rawUrl:         "https://example.com/files/ghjklouiy789iuhjg67t8",
			expectedResult: "ghjklouiy789iuhjg67t8",
		},
		{
			name:           "Get filename by url with empty last segment",
			rawUrl:         "https://example.com/files/",
			expectedResult: "index.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileUrl, _ := url.Parse(tt.rawUrl)

			result := utils.GetFileNameByUrl(fileUrl)
			if result != tt.expectedResult {
				t.Errorf("got %s, expected %s", result, tt.expectedResult)
			}
		})
	}
}

func TestDownloadContent(t *testing.T) {
	tests := []struct {
		name   string
		rawUrl string
		hasErr bool
	}{
		{
			name:   "Download content from simple url",
			rawUrl: "https://static-basket-01.wb.ru/vol1/crm-bnrs/bners1/30_new_big_shopping_day_1440.webp",
			hasErr: false,
		},
		{
			name:   "Download content from url without filename inside",
			rawUrl: "https://www.wildberries.ru/catalog/detyam",
			hasErr: false,
		},
		{
			name:   "Download content from url with empty last segment",
			rawUrl: "https://www.wildberries.ru/",
			hasErr: false,
		},
		{
			name:   "Download content from invalid url",
			rawUrl: "https://www.wild-raspberries.ru/",
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileUrl, _ := url.Parse(tt.rawUrl)
			fileName := utils.GetFileNameByUrl(fileUrl)

			err := utils.DownloadContent(tt.rawUrl, fileName)
			if (err != nil) != tt.hasErr {
				t.Errorf("error = %v, wantErr %v", err, tt.hasErr)
				return
			}

			ok := checkExistFileInCurDir(fileName)

			if ok {
				err := os.Remove(fileName)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if !ok && !tt.hasErr {
				t.Errorf("%v wasn't downloaded from %v", fileName, tt.rawUrl)
			}
		})
	}
}

func TestWriteData(t *testing.T) {
	t.Run("Write data", func(t *testing.T) {
		var output strings.Builder
		data := []string{"line1", "line2", "line3"}

		err := utils.WriteData(&output, data...)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := "line1\nline2\nline3\n"

		if output.String() != expected {
			t.Errorf("got %q, want %q", output.String(), expected)
		}
	})
}

func TestWgetArgs_Parse(t *testing.T) {
	validRawUrl1 := "https://www.wildberries.ru/catalog/detyam"
	validRawUrl2 := "https://www.wildberries.ru/"
	invalidRawUrl := "sdfes://udfsr:sdfsdd:ewf324/234?asfe="

	tests := []struct {
		name     string
		args     []string
		wgetArgs WgetArgs
		hasErr   bool
	}{
		{
			name:     "No args",
			args:     []string{},
			wgetArgs: WgetArgs{},
			hasErr:   true,
		},
		{
			name:     "Invalid url",
			args:     []string{invalidRawUrl},
			wgetArgs: WgetArgs{},
			hasErr:   true,
		},
		{
			name:     "Valid and invalid urls",
			args:     []string{validRawUrl1, invalidRawUrl, validRawUrl2},
			wgetArgs: WgetArgs{},
			hasErr:   true,
		},
		{
			name: "Valid urls",
			args: []string{"https://www.wildberries.ru/catalog/detyam", "https://www.wildberries.ru/"},
			wgetArgs: WgetArgs{
				urls: urlMustParse(validRawUrl1, validRawUrl2),
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			wa := WgetArgs{}

			resetArgs(tt.args)

			err := wa.Parse()

			if (err != nil) != tt.hasErr {
				t.Errorf("WgetArgs.Parse() error = %v, wantErr %v", err, tt.hasErr)
				return
			}

			if len(wa.urls) != len(tt.wgetArgs.urls) {
				t.Errorf("WgetArgs.Parse() urls length = %v, want %v", len(wa.urls), len(tt.wgetArgs.urls))
			}
		})
	}
}

func TestWgetClient_Wget(t *testing.T) {
	name := "Check number errors"
	args := []string{
		"https://static-basket-01.wb.ru/vol1/crm-bnrs/bners1/30_new_big_shopping_day_1440.webp",
		"https://www.wildberries.ru/catalog/detyam",
		"https://www.wildberries.ru/",
		"https://www.wild-raspberries.ru/",
	}
	fileNames := []string{
		"30_new_big_shopping_day_1440.webp",
		"detyam",
		"index.html",
	}
	expectedLenErrs := 1

	t.Run(name, func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(name, flag.ContinueOnError)
		resetArgs(args)

		wc := WgetClient{}

		err := wc.args.Parse()
		if err != nil {
			t.Errorf("not expected error: %q", err)
		}

		errs := wc.Wget()
		lenErrs := len(errs)

		if lenErrs != expectedLenErrs {
			t.Errorf("got %v, expected %v", lenErrs, expectedLenErrs)
		}
	})

	for _, fileName := range fileNames {
		ok := checkExistFileInCurDir(fileName)

		if ok {
			err := os.Remove(fileName)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}
