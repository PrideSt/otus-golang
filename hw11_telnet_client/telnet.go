package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		log.Fatalf("invalid address %q, %s", address, err)
	}

	return &TCPTelnetClient{
		address: tcpAddr,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type TCPTelnetClient struct {
	address net.Addr
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	con     net.Conn
}

func (c *TCPTelnetClient) Connect() error {
	con, err := net.DialTimeout("tcp", c.address.String(), c.timeout)
	if err != nil {
		return err
	}
	c.con = con

	return nil
}

func (c *TCPTelnetClient) Close() error {
	if err := c.in.Close(); err != nil {
		return fmt.Errorf("close connection fault: %w", err)
	}

	if err := c.con.Close(); err != nil {
		return fmt.Errorf("close connection fault: %w", err)
	}

	return nil
}

func (c *TCPTelnetClient) Send() error {
	_, err := io.Copy(c.con, c.in)
	return err
}

func (c *TCPTelnetClient) Receive() error {
	_, err := io.Copy(c.out, c.con)
	return err
}
