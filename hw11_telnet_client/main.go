package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	timeout time.Duration
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	flag.Parse()
	host := flag.Arg(0)
	port := flag.Arg(1)

	if len(flag.Args()) < 2 {
		log.Fatal("invalid arguments, host and port required. Usage: go-telnet [--timeout=10s] <host> <port>")
	}

	log.SetOutput(os.Stderr)
	log.Println("pid:", os.Getpid())
	log.Println("timeout:", timeout)
	log.Println("host:", host)
	log.Println("port:", port)

	addr := net.JoinHostPort(host, port)
	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	chSigTerm := handleSignals(ctx)

	chSenderTerm := runSender(client)
	chReceiverTerm := runReceiver(client)

	select {
	case <-chSigTerm:
		log.Println("signal handled, terminate")
	case <-chSenderTerm:
		log.Println("stdin closed, terminate")
		cancel()
	case <-chReceiverTerm:
		log.Println("connection closed by peer, terminate")
		cancel()
	}

	log.Println("close telnet client")
	if err := client.Close(); err != nil {
		log.Fatal(err)
	}

	<-chSigTerm
	<-chReceiverTerm
	// don't waiting for stdin closing
	// https://stackoverflow.com/questions/63789503/is-it-possible-to-interrupt-io-copy-by-closing-src
	// <-chSenderTerm
	log.Println("bye!")
}

func handleSignals(ctx context.Context) <-chan struct{} {
	chTerm := make(chan struct{})

	go func() {
		defer close(chTerm)

		sigChan := make(chan os.Signal, 1)
		defer close(sigChan)

		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigChan)

		select {
		case sign := <-sigChan:
			log.Printf("signal %s recv, terminate...", sign)

			return

		case <-ctx.Done():
			log.Println("signal handler terminated by context")

			return
		}
	}()

	return chTerm
}

func runSender(client TelnetClient) <-chan struct{} {
	chSenderTerm := make(chan struct{})
	go func() {
		defer close(chSenderTerm)
		if err := client.Send(); err != nil {
			log.Println("Send error:", err)

			return
		}
		log.Println("Sender EOF!")
	}()

	return chSenderTerm
}

func runReceiver(client TelnetClient) <-chan struct{} {
	chReceiverTerm := make(chan struct{})
	go func() {
		defer close(chReceiverTerm)
		if err := client.Receive(); err != nil {
			log.Println("Recv error:", err)

			return
		}
		log.Println("Connection closed by peer")
	}()

	return chReceiverTerm
}
