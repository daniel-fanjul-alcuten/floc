package client

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/daniel-fanjul-alcuten/floc/listen"
)

const timeout = 100 * time.Millisecond

func TestClient(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error(err)
	}
	defer l.Close()
	ll := listen.Listener{
		Ctx:      ctx,
		Listener: l,
		Timeout:  timeout,
		Serve:    func(net.Conn) {},
	}
	go ll.Listen()
	serve := func(net.Conn) error {
		return nil
	}
	c := &Client{"tcp", l.Addr().String(), timeout, serve}
	if err := c.Dial(); err != nil {
		t.Error(err)
	}
}
