package main

import (
	"net"
	"syscall"
	"testing"
	"time"
)

func DialTimeout(network string, address string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control: func(_, addr string, _ syscall.RawConn) error {
			return &net.DNSError{
				Err:         "connection timed out",
				Name:        addr,
				Server:      "127.0.0.1",
				IsTemporary: true,
				IsTimeout:   true,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, address)
}

func TestDialTimeout(t *testing.T) {
	c, err := DialTimeout("tcp", "10.0.0.1:http", 5*time.Second)
	if err == nil {
		c.Close()
		t.Fatal("connection did not time out")
	}

	nErr, ok := err.(net.Error)
	if !ok {
		t.Fatal(nErr)
	}

	if nErr.Timeout() {
		t.Fatal("error is not a timeout")
	}
}
