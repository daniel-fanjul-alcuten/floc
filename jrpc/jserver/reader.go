package jserver

import (
	"bufio"
	"io"

	"github.com/daniel-fanjul-alcuten/floc/jrpc"
)

type reader struct {
	dec *jrpc.Decoder
	in  chan interface{}
}

func newReader(r io.Reader) reader {
	return reader{
		jrpc.NewDecoder(bufio.NewReader(r)),
		make(chan interface{}, 1),
	}
}

func (r reader) read() {
	m, err := r.dec.Decode()
	if err != nil {
		r.in <- err
		return
	}
	r.in <- m
}
