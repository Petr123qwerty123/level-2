package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные из сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

const (
	protocolType   = "tcp"
	defaultTimeout = 10 * time.Second
)

var errNArgs = errors.New("invalid number of arguments")

// TelnetFlags структура, определяющюая опции утилиты Telnet
type TelnetFlags struct {
	timeout time.Duration
}

// Parse метод для распарсивания и сохранения значений флагов опций в поля структуры TelnetFlags
func (tf *TelnetFlags) Parse() {
	flag.DurationVar(&tf.timeout, "timeout", defaultTimeout, "Specify timeout")

	flag.Parse()
}

// TelnetArgs структура, определяющая неименованные аргументы запуска утилиты Telnet
type TelnetArgs struct {
	host string
	port string
}

// Parse метод для распарсивания и сохранения значений неименованных аргументов запуска утилиты Telnet в поля структуры
// TelnetArgs
func (ta *TelnetArgs) Parse() error {
	nArg := flag.NArg()

	switch nArg {
	case 2:
		ta.host = flag.Arg(0)
		ta.port = flag.Arg(1)
	default:
		return errNArgs
	}

	return nil
}

// TelnetClient структура для управления утилитой Telnet
type TelnetClient struct {
	conn    net.Conn
	flags   TelnetFlags
	args    TelnetArgs
	stopSig chan os.Signal
	sync.WaitGroup
}

// NewTelnetClient конструктор для создания объекта структуры TelnetClient
func NewTelnetClient() (*TelnetClient, error) {
	// создаем канал, ожидающий сигнал о завершении работы утилиты
	stopSig := make(chan os.Signal, 1)

	tc := &TelnetClient{stopSig: stopSig}

	tc.flags.Parse()
	err := tc.args.Parse()
	if err != nil {
		return nil, err
	}

	return tc, nil
}

// Connect метод для подключения к tcp-серверу
func (tc *TelnetClient) Connect() error {
	address := tc.args.host + ":" + tc.args.port

	conn, err := net.DialTimeout(protocolType, address, tc.flags.timeout)
	if err != nil {
		return err
	}

	tc.conn = conn

	return nil
}

// receiveMessages метод, который переносит поток данных из сокета в os.Stdout, в случае возникновения ошибки в канал,
// ожидающий сигнала о завершении работы утилиты, передается os.Kill, если он не закрыт
func (tc *TelnetClient) receiveMessages() {
	defer tc.Done()

	_, err := io.Copy(os.Stdout, tc.conn)
	if err != nil {
		if _, ok := <-tc.stopSig; ok {
			tc.stopSig <- os.Kill
		}
	}
}

// sendMessages метод, который переносит поток данных из os.Stdin в сокет, в случае возникновения ошибки в канал,
// ожидающий сигнала о завершении работы утилиты, передается os.Kill, если он не закрыт
func (tc *TelnetClient) sendMessages() {
	defer tc.Done()

	_, err := io.Copy(tc.conn, os.Stdin)
	if err != nil {
		if _, ok := <-tc.stopSig; ok {
			tc.stopSig <- os.Kill
		}
	}
}

// Start метод запуска утилиты
func (tc *TelnetClient) Start() error {
	// подключаемся к tcp серверу, в случае возникновения ошибки закрываем канал, ожидающий сигнала о завершении работы
	// утилиты, закрываем подключение к tcp-серверу
	err := tc.Connect()
	if err != nil {
		close(tc.stopSig)
		err = tc.conn.Close()
		return err
	}

	tc.Add(2)

	// запускаем горутины, связанные с передачей сообщений из os.Stdin в сокет, из сокета в os.Stdout
	go tc.receiveMessages()
	go tc.sendMessages()

	return nil
}

// Stop метод для остановки работы утилиты
// дожидается завершения работы горутин, связанных с передачей сообщений из os.Stdin в сокет, из сокета в os.Stdout,
// закрываем канал, ожидающий сигнала о завершении работы утилиты, закрываем подключение к tcp-серверу
func (tc *TelnetClient) Stop() error {
	defer tc.Wait()

	close(tc.stopSig)

	err := tc.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// создание объекта структуры TelnetClient
	telnet, err := NewTelnetClient()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}

	// запуск утилиты Telnet, в случае ошибки - её вывод и выход из программы с кодом ошибки 1
	err = telnet.Start()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}

	// отслеживание сигнала о завершении работы утилиты, полученного от пользователя и его запись в канал telnet.stopSig
	signal.Notify(telnet.stopSig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	// ожидание получения сигнала о завершении работы утилиты, полученного от пользователя
	<-telnet.stopSig

	// остановка работы утилиты, в случае ошибки - её вывод и выход из программы с кодом ошибки 1
	err = telnet.Stop()
	if err != nil {
		fmt.Printf("%q\n", err)
		os.Exit(1)
	}

	fmt.Println("Получен сигнал завершения программы")
}
