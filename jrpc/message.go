// Package jrpc implements http://www.jsonrpc.org/specification_v1
package jrpc

import (
	"encoding/json"
	"io"
)

// Request is an JSON RPC v1 request message.
type Request struct {
	ID     interface{}
	Method string
	Params []interface{}
}

// Response is an JSON RPC v1 response message.
type Response struct {
	ID     interface{}
	Result interface{}
	Error  interface{}
}

// Notification is an JSON RPC v1 notification message.
type Notification struct {
	Method string
	Params []interface{}
}

// Encoder is a json.Encoder with custom methods to properly encode these types.
type Encoder json.Encoder

// NewEncoder returns a new Encoder that does not escape HTML.
func NewEncoder(w io.Writer) *Encoder {
	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)
	e.SetIndent("", "")
	return (*Encoder)(e)
}

type request struct {
	ID     interface{}   `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type response struct {
	ID     interface{} `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

// EncodeRequest encodes a Request.
func (e *Encoder) EncodeRequest(r Request) error {
	return (*json.Encoder)(e).Encode(request{
		ID:     r.ID,
		Method: r.Method,
		Params: r.Params,
	})
}

// EncodeResponse encodes a Response. Only the Error or the Result is written.
func (e *Encoder) EncodeResponse(r Response) error {
	if r.Error != nil {
		return (*json.Encoder)(e).Encode(response{
			ID:    r.ID,
			Error: r.Error,
		})
	}
	return (*json.Encoder)(e).Encode(response{
		ID:     r.ID,
		Result: r.Result,
	})
}

// EncodeNotification encodes a Notification.
func (e *Encoder) EncodeNotification(n Notification) error {
	return (*json.Encoder)(e).Encode(request{
		Method: n.Method,
		Params: n.Params,
	})
}

// Decoder is a json.Decoder with custom methods to properly decode these types.
type Decoder json.Decoder

// NewDecoder returns a new Decoder that disallows unknown fields.
func NewDecoder(r io.Reader) *Decoder {
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	return (*Decoder)(d)
}

type message struct {
	ID     interface{}   `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	Result interface{}   `json:"result"`
	Error  interface{}   `json:"error"`
}

// Decode decodes a Request, Response or Notification.
func (d *Decoder) Decode() (m interface{}, err error) {
	var j message
	if err = (*json.Decoder)(d).Decode(&j); err != nil {
		return
	}
	if j.Method != "" {
		if j.ID != nil {
			m = Request{
				ID:     j.ID,
				Method: j.Method,
				Params: j.Params,
			}
			return
		}
		m = Notification{
			Method: j.Method,
			Params: j.Params,
		}
		return
	}
	if j.Error != nil {
		m = Response{
			ID:    j.ID,
			Error: j.Error,
		}
		return
	}
	m = Response{
		ID:     j.ID,
		Result: j.Result,
	}
	return
}
