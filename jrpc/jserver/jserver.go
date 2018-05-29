// Package jserver implements a Server using the package jrpc
package jserver

import (
	"bytes"
	"context"
	"io"
	"math/big"

	"github.com/daniel-fanjul-alcuten/floc/jrpc"
)

// Request is a pair method and params that expects a Response.
type Request struct {
	Method string
	Params []interface{}
}

// Response is a pair result or error.
type Response struct {
	Result interface{}
	Error  interface{}
}

// Notification is a pair method and params that does not expect a
// Response.
type Notification struct {
	Method string
	Params []interface{}
}

// RequestResponse is a Request with a channel to set the Response.
type RequestResponse struct {
	Request
	Response chan<- Response
}

// idResponse is the Response that has been returned for the Request with the
// given id.
type idResponse struct {
	id       interface{}
	response Response
}

// Server sends and receives Requests, Responses and Notifications through an
// io.Reader and an io.Writer.
type Server struct {
	ctx               context.Context
	cancel            context.CancelFunc
	reader            reader
	writer            writer
	recvRequests      chan RequestResponse
	sendRequests      chan RequestResponse
	recvNotifications chan Notification
	sendNotifications chan Notification
	id                *big.Int
	responses         map[string]chan<- Response
	idResponses       chan idResponse
}

// NewServer returns a new Server that reads from r, writes to w, uses channels
// of size recvSize to receive and uses channels of size sendSize to send.
func NewServer(ctx context.Context, r io.Reader, w io.Writer, recvSize, sendSize int) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:               ctx,
		cancel:            cancel,
		reader:            newReader(r),
		writer:            newWriter(w),
		recvRequests:      make(chan RequestResponse, recvSize),
		sendRequests:      make(chan RequestResponse, sendSize),
		recvNotifications: make(chan Notification, recvSize),
		sendNotifications: make(chan Notification, sendSize),
		id:                big.NewInt(0),
		responses:         make(map[string]chan<- Response, 32),
		idResponses:       make(chan idResponse),
	}
}

// RecvRequests returns a channel where RequestResponses are received.
func (s *Server) RecvRequests() <-chan RequestResponse {
	return s.recvRequests
}

// SendRequest sends a Request, waits until the Response arrives and returns
// it.
func (s *Server) SendRequest(r Request) Response {
	ch := make(chan Response, 1)
	s.sendRequests <- RequestResponse{r, ch}
	return <-ch
}

// RecvNotifications returns a channel where Notifications are received.
func (s *Server) RecvNotifications() <-chan Notification {
	return s.recvNotifications
}

// SendNotification sends a Notification.
func (s *Server) SendNotification(n Notification) {
	s.sendNotifications <- n
}

var one = big.NewInt(1)

// Serve reads and writes messages until there is an error or the context is
// done.
func (s *Server) Serve() error {
	defer s.cancel()
	go s.reader.read()
	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case r := <-s.sendRequests:
			s.sendRequest(r)
		case n := <-s.sendNotifications:
			s.sendNotification(n)
		case err := <-s.writer.out:
			if err != nil {
				return err
			}
			s.writer.written()
		case m := <-s.reader.in:
			switch m := m.(type) {
			case jrpc.Request:
				s.recvRequest(m)
			case jrpc.Notification:
				s.recvNotification(m)
			case jrpc.Response:
				s.recvResponse(m)
			case error:
				return m
			}
			go s.reader.read()
		case idResp := <-s.idResponses:
			s.sendResponse(idResp)
		}
	}
}

func (s *Server) sendRequest(r RequestResponse) {
	id := s.id.Add(s.id, one).Text(62)
	buf := &bytes.Buffer{}
	enc := jrpc.NewEncoder(buf)
	if err := enc.EncodeRequest(jrpc.Request{
		ID:     id,
		Method: r.Method,
		Params: r.Params,
	}); err == nil {
		s.writer.add(buf)
		s.responses[id] = r.Response
	}
}

func (s *Server) sendNotification(n Notification) {
	buf := &bytes.Buffer{}
	enc := jrpc.NewEncoder(buf)
	if err := enc.EncodeNotification(jrpc.Notification{
		Method: n.Method,
		Params: n.Params,
	}); err == nil {
		s.writer.add(buf)
	}
}

func (s *Server) recvRequest(r jrpc.Request) {
	ch := make(chan Response, 1)
	s.recvRequests <- RequestResponse{Request{r.Method, r.Params}, ch}
	go func() {
		select {
		case <-s.ctx.Done():
		case s.idResponses <- idResponse{r.ID, <-ch}:
		}
	}()
}

func (s *Server) recvNotification(n jrpc.Notification) {
	s.recvNotifications <- Notification{n.Method, n.Params}
}

func (s *Server) recvResponse(r jrpc.Response) {
	id, ok := r.ID.(string)
	if !ok {
		return
	}
	ch := s.responses[id]
	if ch == nil {
		return
	}
	ch <- Response{r.Result, r.Error}
	delete(s.responses, id)
}

func (s *Server) sendResponse(idResp idResponse) {
	buf := &bytes.Buffer{}
	enc := jrpc.NewEncoder(buf)
	if err := enc.EncodeResponse(jrpc.Response{
		ID:     idResp.id,
		Result: idResp.response.Result,
		Error:  idResp.response.Error,
	}); err == nil {
		s.writer.add(buf)
	}
}
