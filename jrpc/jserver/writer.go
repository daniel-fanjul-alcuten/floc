package jserver

import (
	"bufio"
	"bytes"
	"io"
)

type writer struct {
	w       *bufio.Writer
	bufs    []*bytes.Buffer
	writing bool
	out     chan error
}

func newWriter(w io.Writer) writer {
	return writer{
		bufio.NewWriter(w),
		nil,
		false,
		make(chan error, 1),
	}
}

func (w writer) add(buf *bytes.Buffer) {
	w.bufs = append(w.bufs, buf)
	if !w.writing {
		go w.write(w.bufs)
		w.writing, w.bufs = true, nil
	}
}

func (w writer) write(bufs []*bytes.Buffer) {
	for _, buf := range bufs {
		if _, err := w.w.Write(buf.Bytes()); err != nil {
			w.out <- err
			return
		}
	}
	w.out <- w.w.Flush()
}

func (w writer) written() {
	if len(w.bufs) > 0 {
		go w.write(w.bufs)
		w.bufs = nil
	} else {
		w.writing = false
	}
}
