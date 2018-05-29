package listen

import (
	"context"
	"net"
	"testing"
	"time"
)

const timeout = 100 * time.Millisecond

func TestListen_NoConn(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error(err)
	}
	defer l.Close()
	ll := Listener{ctx, l, timeout, nil}
	if err := ll.Listen(); err != context.DeadlineExceeded {
		t.Error(err)
	}
}

func TestListen_TCPConn(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error(err)
	}
	defer l.Close()
	go func() {
		conn, err := net.DialTimeout("tcp", l.Addr().String(), timeout)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
	}()
	conns := make(chan net.Conn, 1)
	serve := func(conn net.Conn) {
		conns <- conn
	}
	ll := Listener{ctx, l, timeout, serve}
	if err := ll.Listen(); err != context.DeadlineExceeded {
		t.Error(err)
	}
	select {
	case <-conns:
	case <-time.After(timeout):
		t.Error()
	}
}

func TestListen_UnixConn(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	l, err := net.Listen("unix", "./test.socket")
	if err != nil {
		t.Error(err)
	}
	defer l.Close()
	go func() {
		conn, err := net.DialTimeout("unix", "./test.socket", timeout)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
	}()
	conns := make(chan net.Conn, 1)
	serve := func(conn net.Conn) {
		conns <- conn
	}
	ll := Listener{ctx, l, timeout, serve}
	if err := ll.Listen(); err != context.DeadlineExceeded {
		t.Error(err)
	}
	select {
	case <-conns:
	case <-time.After(timeout):
		t.Error()
	}
}
