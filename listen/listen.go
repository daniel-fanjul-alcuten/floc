package listen

import (
	"context"
	"net"
	"time"
)

// Listener is a wrapper to net.Listener.
type Listener struct {
	Ctx      context.Context
	Listener net.Listener
	Timeout  time.Duration
	Serve    func(net.Conn)
}

type deadlineSettable interface {
	SetDeadline(t time.Time) error
}

// Listen calls l.Listener.Accept() in a loop until the context l.Ctx is done
// or there is an error. time.Now().Add(l.Timeout) is passed to
// l.Listener.SetDeadline() before each call to Accept(), and errors caused by
// this deadline are ignored. Every accepted connection is passed to the
// callback l.Serve() in a new goroutine that closes it afterwards.
func (l *Listener) Listen() error {
	ds := l.Listener.(deadlineSettable)
	for {
		select {
		case <-l.Ctx.Done():
			return l.Ctx.Err()
		default:
			err := ds.SetDeadline(time.Now().Add(l.Timeout))
			if err != nil {
				return err
			}
			conn, err := l.Listener.Accept()
			if err != nil {
				if err.(net.Error).Timeout() || err.(net.Error).Temporary() {
					continue
				}
				return err
			}
			go func() {
				defer conn.Close()
				l.Serve(conn)
			}()
		}
	}
}
