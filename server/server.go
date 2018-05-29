package server

import (
	"context"
	"net"
	"time"

	"github.com/daniel-fanjul-alcuten/floc/listen"
)

// Server announces on an address and invokes a callback to handle the
// connections.
type Server struct {
	Ctx     context.Context
	Network string
	Address string
	Timeout time.Duration
	Serve   func(net.Conn)
}

// Listen announces on s.Network and s.Address and calls and returns
// listen.Listener.Listen().
func (s *Server) Listen() error {
	l, err := net.Listen(s.Network, s.Address)
	if err != nil {
		return err
	}
	defer l.Close()
	ll := &listen.Listener{
		Ctx:      s.Ctx,
		Listener: l,
		Timeout:  s.Timeout,
		Serve:    s.Serve,
	}
	return ll.Listen()
}
