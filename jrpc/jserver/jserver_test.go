package jserver

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

const timeout = 100 * time.Millisecond

func TestServe_Trivial(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	reader := &bytes.Buffer{}
	writer := &bytes.Buffer{}
	server := NewServer(ctx, reader, writer, 32, 32)
	if err := server.Serve(); err != io.EOF {
		t.Error(err)
	}
}

func TestServe_SendRequest(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rreader, rwriter := io.Pipe()
	wreader, wwriter := io.Pipe()
	server := NewServer(ctx, rreader, wwriter, 32, 32)
	go func() {
		server.SendRequest(Request{"echo", []interface{}{1, "2"}})
		cancel()
	}()
	go readUntilAndWrite(t,
		wreader, `{"id":"1","method":"echo","params":[1,"2"]}`,
		rwriter, `{"id":"1","result":"ok"}`,
	)
	if err := server.Serve(); err != context.Canceled {
		t.Error(err)
	}
}

func TestServe_SendRequest_FailedWrite(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rreader, _ := io.Pipe()
	writer := &bytes.Buffer{}
	server := NewServer(ctx, rreader, &TruncatingWriter{writer, 10}, 32, 32)
	go func() {
		server.SendRequest(Request{"echo", []interface{}{1, "2"}})
	}()
	if err := server.Serve(); err != errTruncated {
		t.Error(err)
	}
}

func TestServe_SendNotification(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rreader, _ := io.Pipe()
	wreader, wwriter := io.Pipe()
	server := NewServer(ctx, rreader, wwriter, 32, 32)
	go func() {
		server.SendNotification(Notification{"echo", []interface{}{1, "2"}})
	}()
	go readUntilAndCancel(t,
		wreader, `{"id":null,"method":"echo","params":[1,"2"]}`,
		cancel,
	)
	if err := server.Serve(); err != context.Canceled {
		t.Error(err)
	}
}

func TestServe_RecvRequest_Response(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rreader, rwriter := io.Pipe()
	wreader, wwriter := io.Pipe()
	server := NewServer(ctx, rreader, wwriter, 32, 32)
	go func() {
		if _, err := rwriter.Write([]byte(`{"id":0,"method":"method1","params":null}`)); err != nil {
			t.Error(err)
		}
	}()
	go func() {
		r := <-server.RecvRequests()
		r.Response <- Response{"ok", nil}
	}()
	go readUntilAndCancel(t,
		wreader, `{"id":0,"result":"ok","error":null}`,
		cancel,
	)
	if err := server.Serve(); err != context.Canceled {
		t.Error(err)
	}
}

func TestServe_RecvNotification(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rreader, rwriter := io.Pipe()
	writer := &bytes.Buffer{}
	server := NewServer(ctx, rreader, writer, 32, 32)
	go func() {
		if _, err := rwriter.Write([]byte(`{"id":null,"method":"method1","params":null}`)); err != nil {
			t.Error(err)
		}
	}()
	go func() {
		<-server.RecvNotifications()
		cancel()
	}()
	if err := server.Serve(); err != context.Canceled {
		t.Error(err)
	}
}

func TestServe_RecvResponse_WrongID_Type(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rreader, rwriter := io.Pipe()
	writer := &bytes.Buffer{}
	server := NewServer(ctx, rreader, writer, 32, 32)
	go func() {
		if _, err := rwriter.Write([]byte(`{"id":1,"result":"ok"}`)); err != nil {
			t.Error(err)
		}
	}()
	if err := server.Serve(); err != context.DeadlineExceeded {
		t.Error(err)
	}
}

func TestServe_RecvResponse_WrongID_Value(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rreader, rwriter := io.Pipe()
	writer := &bytes.Buffer{}
	server := NewServer(ctx, rreader, writer, 32, 32)
	go func() {
		if _, err := rwriter.Write([]byte(`{"id":"1","error":"ko"}`)); err != nil {
			t.Error(err)
		}
	}()
	if err := server.Serve(); err != context.DeadlineExceeded {
		t.Error(err)
	}
}

func readUntilAndWrite(t *testing.T, r io.Reader, input string, w io.Writer, output string) {
	buf := &bytes.Buffer{}
	for {
		if _, err := io.CopyN(buf, r, 1); err != nil {
			t.Error(err)
		}
		if s := buf.String(); s == input+"\n" {
			if _, err := w.Write([]byte(output)); err != nil {
				t.Error(err)
			}
			break
		} else if !strings.HasPrefix(input, buf.String()) {
			t.Error(buf.String())
		}
	}
}

func readUntilAndCancel(t *testing.T, r io.Reader, input string, cancel context.CancelFunc) {
	buf := &bytes.Buffer{}
	for {
		if _, err := io.CopyN(buf, r, 1); err != nil {
			t.Error(err)
		}
		if s := buf.String(); s == input+"\n" {
			cancel()
			break
		} else if !strings.HasPrefix(input, buf.String()) {
			t.Error(buf.String())
		}
	}
}
