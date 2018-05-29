package client

import (
	"net"
	"time"
)

// Client dials a connection and invokes a callback to handle it.
type Client struct {
	Network string
	Address string
	Timeout time.Duration
	Serve   func(net.Conn) error
}

// Dial dials a connection to the c.Network and c.Address with the timeout
// c.Timeout, calls f.Serve, waits for it to finish, closes the connection and
// returns any error.
func (c *Client) Dial() (err error) {
	conn, err := net.DialTimeout(c.Network, c.Address, c.Timeout)
	if err != nil {
		return
	}
	defer conn.Close()
	return c.Serve(conn)
}
