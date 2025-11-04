package main

import (
	"bufio"
	"context"
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

func main() {
	var host string
	var port string
	var timeout int

	dst := os.Stdout
	src := os.Stdin
	stdErr := os.Stderr

	flag.StringVar(&host, "h", "", "host")
	flag.StringVar(&port, "p", "", "port")
	flag.IntVar(&timeout, "t", 10, "timeout in seconds")
	flag.Parse()

	if host == "" || port == "" {
		fmt.Fprintln(stdErr, "use: telnet --host <host> --port <port> [--timeout <seconds>]")
		os.Exit(1)
	}

	address := net.JoinHostPort(host, port)
	dialer := net.Dialer{
		Timeout: time.Duration(timeout) * time.Second,
	}
	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		fmt.Fprintf(stdErr, "dial error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Fprintf(stdErr, "connected to %s\n", address)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	wg.Go(func() {
		defer cancel()
		reader := bufio.NewReader(conn)
		buf := make([]byte, 1<<9) // 1024
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			n, err := reader.Read(buf)
			if err != nil {
				var netErr net.Error
				if errors.As(err, &netErr) && netErr.Timeout() {
					continue
				}
				if err == io.EOF {
					fmt.Fprintln(stdErr, "connection closed by server")
					return
				}
				if err != io.EOF && ctx.Err() == nil {
					fmt.Fprintf(stdErr, "error reading from socket: %v\n", err)
				}
				return
			}
			if n > 0 {
				fmt.Fprint(dst, string(buf[:n]))
			}

		}
	})
	wg.Go(func() {
		defer cancel()
		reader := bufio.NewReader(src)
		writer := bufio.NewWriter(conn)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Fprintln(stdErr, "closing connection")
					return
				}
				if ctx.Err() == nil {
					fmt.Fprintf(stdErr, "stdin read error: %v\n", err)
				}
				return
			}

			_, err = writer.WriteString(line)
			if err != nil {
				fmt.Fprintf(stdErr, "write error: %v\n", err)
				return
			}
			err = writer.Flush()
			if err != nil {
				fmt.Fprintf(stdErr, "flush error: %v\n", err)
				return
			}
		}
	})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	go func() {
		<-sigCh
		fmt.Fprintln(dst, "\ninterrupt signal received")
		cancel()
	}()
	<-ctx.Done()

}
