package main

import (
	"context"
	"flag"
	"net"
	"os"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

type tcpServer struct {
	host    string
	port    uint16
	ctx     context.Context
	ln      net.Listener
	conns   []net.Conn
	stopSig chan os.Signal
	sync.WaitGroup
}

func newTcpServer(port uint16, ctx context.Context) *tcpServer {
	stopSig := make(chan os.Signal, 1)

	ts := &tcpServer{
		host:    "localhost",
		port:    port,
		stopSig: stopSig,
		ctx:     ctx,
	}

	return ts
}

func (ts *tcpServer) listenPort() error {
	var lc net.ListenConfig

	ln, err := lc.Listen(ts.ctx, protocolType, ts.host+":"+strconv.Itoa(int(ts.port)))
	if err != nil {
		return err
	}

	ts.ln = ln

	return nil
}

func (ts *tcpServer) acceptConnections() {
	defer ts.Done()

	for {
		conn, err := ts.ln.Accept()
		if err != nil {
			if _, ok := <-ts.stopSig; ok {
				ts.stopSig <- os.Kill
			}

			return
		}

		ts.conns = append(ts.conns, conn)
	}
}

func (ts *tcpServer) start() error {
	err := ts.listenPort()
	if err != nil {
		close(ts.stopSig)
		err = ts.ln.Close()
		return err
	}

	ts.Add(1)

	go ts.acceptConnections()

	return nil
}

func (ts *tcpServer) stop() {
	defer ts.Wait()

	close(ts.stopSig)

	for _, conn := range ts.conns {
		conn.Close()
	}

	ts.ln.Close()
}

// Helper function to reset the command-line args.
func resetArgs(args []string) {
	os.Args = []string{"testArgs"}

	for _, arg := range args {
		os.Args = append(os.Args, arg)
	}
}

func TestTelnetFlags_Parse(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		flags TelnetFlags
	}{
		{
			name: "No flags",
			args: []string{},
			flags: TelnetFlags{
				timeout: 10 * time.Second,
			},
		},
		{
			name: "Timeout flag",
			args: []string{"--timeout", "20s"},
			flags: TelnetFlags{
				timeout: 20 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			tf := TelnetFlags{}

			resetArgs(tt.args)
			tf.Parse()

			if !reflect.DeepEqual(tf, tt.flags) {
				t.Errorf("TelnetFlags.Parse() got = %v, want %v", tf, tt.flags)
			}
		})
	}
}

func TestTelnetArgs_Parse(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		telnetArgs TelnetArgs
		hasErr     bool
	}{
		{
			name:       "Invalid number of arguments",
			args:       []string{},
			telnetArgs: TelnetArgs{},
			hasErr:     true,
		},
		{
			name: "Valid number of arguments",
			args: []string{"sdfsdfs", "asdasdad"},
			telnetArgs: TelnetArgs{
				host: "sdfsdfs",
				port: "asdasdad",
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			ta := TelnetArgs{}

			resetArgs(tt.args)

			flag.Parse()
			err := ta.Parse()

			if (err != nil) != tt.hasErr {
				t.Errorf("TelnetArgs.Parse() error = %v, wantErr %v", err, tt.hasErr)
				return
			}

			if !reflect.DeepEqual(tt.telnetArgs, ta) {
				t.Errorf("TelnetArgs.Parse() got %v, want %v", ta, tt.telnetArgs)
			}
		})
	}
}

func TestTelnetClient_Connect(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		conn   net.Conn
		hasErr bool
	}{
		{
			name:   "Server is not waiting for a connection (or not enough data to connect)",
			args:   []string{"localhost", "3242"},
			hasErr: true,
		},
		{
			name:   "Valid host and port",
			args:   []string{"localhost", "1234"},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			resetArgs(tt.args)

			tc := TelnetClient{}

			tc.flags.Parse()
			err = tc.args.Parse()
			if err != nil {
				t.Errorf("not expected error: %q", err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), tc.flags.timeout)
			defer cancel()

			ts := newTcpServer(1234, ctx)
			if err != nil {
				t.Errorf("not expected error: %q", err)
				return
			}

			err = ts.start()
			if err != nil {
				t.Errorf("not expected error: %q", err)
				return
			}

			err = tc.Connect()
			if (err != nil) != tt.hasErr {
				t.Errorf("TelnetClient.Connect() error = %v, wantErr %v", err, tt.hasErr)
				return
			}

			ts.stop()

			result := len(ts.conns)
			expectedResult := 1

			if (result != expectedResult) && !tt.hasErr {
				t.Errorf("got %v, expected %v", result, expectedResult)
			}
		})
	}
}
