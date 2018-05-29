package jserver

import (
	"errors"
	"io"
)

type TruncatingWriter struct {
	w io.Writer
	n int
}

var errTruncated = errors.New("Truncated")

func (w TruncatingWriter) Write(p []byte) (n int, err error) {
	q := p
	if len(q) > w.n {
		q = q[:w.n]
	}
	n, err = w.w.Write(q)
	if err == nil && n < len(p) {
		err = errTruncated
	}
	w.n -= n
	return
}
