package server

import (
	"context"
	"net"
	"testing"
	"time"
)

const timeout = 100 * time.Millisecond

func TestServer(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	s := &Server{ctx, "tcp", "127.0.0.1:0", timeout, func(net.Conn) {}}
	if err := s.Listen(); err != context.DeadlineExceeded {
		t.Error(err)
	}
}
