package jrpc

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEncodeRequest(t *testing.T) {
	type entry struct {
		given Request
		then  string
	}
	for i, e := range []entry{
		// 0
		{Request{0, "method1", nil},
			`{"id":0,"method":"method1","params":null}`},
		// 1
		{Request{"1", "method1", nil},
			`{"id":"1","method":"method1","params":null}`},
		// 2
		{Request{"2", "method1", []interface{}{}},
			`{"id":"2","method":"method1","params":[]}`},
		// 3
		{Request{"3", "method1", []interface{}{1, "2"}},
			`{"id":"3","method":"method1","params":[1,"2"]}`},
	} {
		buf := bytes.Buffer{}
		enc := NewEncoder(&buf)
		if err := enc.EncodeRequest(e.given); err != nil {
			t.Error(err)
		} else if s := buf.String(); s != e.then+"\n" {
			t.Errorf("%v: %#v: %v", i, e.given, s)
		}
	}
}

func TestEncodeResponse(t *testing.T) {
	type entry struct {
		given Response
		then  string
	}
	for i, e := range []entry{
		// 0
		{Response{0, 1, nil},
			`{"id":0,"result":1,"error":null}`},
		// 1
		{Response{"1", "2", nil},
			`{"id":"1","result":"2","error":null}`},
		// 2
		{Response{0, nil, 1},
			`{"id":0,"result":null,"error":1}`},
		// 3
		{Response{"1", nil, "2"},
			`{"id":"1","result":null,"error":"2"}`},
	} {
		buf := bytes.Buffer{}
		enc := NewEncoder(&buf)
		if err := enc.EncodeResponse(e.given); err != nil {
			t.Error(err)
		} else if s := buf.String(); s != e.then+"\n" {
			t.Errorf("%v: %#v: %v", i, e.given, s)
		}
	}
}

func TestEncodeNotification(t *testing.T) {
	type entry struct {
		given Notification
		then  string
	}
	for i, e := range []entry{
		// 0
		{Notification{"method1", nil},
			`{"id":null,"method":"method1","params":null}`},
		// 1
		{Notification{"method1", []interface{}{}},
			`{"id":null,"method":"method1","params":[]}`},
		// 2
		{Notification{"method1", []interface{}{1, "2"}},
			`{"id":null,"method":"method1","params":[1,"2"]}`},
	} {
		buf := bytes.Buffer{}
		enc := NewEncoder(&buf)
		if err := enc.EncodeNotification(e.given); err != nil {
			t.Error(err)
		} else if s := buf.String(); s != e.then+"\n" {
			t.Errorf("%v: %#v: %v", i, e.given, s)
		}
	}
}

func TestDecode(t *testing.T) {
	type entry struct {
		given string
		then  interface{}
		err   string
	}
	for i, e := range []entry{
		// 0
		{`{"id":0,"method":"method1","params":null}`,
			Request{0, "method1", nil},
			""},
		// 1
		{`{"id":"1","method":"method1","params":null}`,
			Request{"1", "method1", nil},
			""},
		// 2
		{`{"id":"2","method":"method1","params":[]}`,
			Request{"2", "method1", []interface{}{}},
			""},
		// 3
		{`{"id":"3","method":"method1","params":[1,"2"]}`,
			Request{"3", "method1", []interface{}{1, "2"}},
			""},
		// 4
		{`{"id":0,"result":1,"error":null}`,
			Response{0, 1, nil},
			""},
		// 5
		{`{"id":"1","result":"2","error":null}`,
			Response{"1", "2", nil},
			""},
		// 6
		{`{"id":0,"result":null,"error":1}`,
			Response{0, nil, 1},
			""},
		// 7
		{`{"id":"1","result":null,"error":"2"}`,
			Response{"1", nil, "2"},
			""},
		// 8
		{`{"id":null,"method":"method1","params":null}`,
			Notification{"method1", nil},
			""},
		// 9
		{`{"id":null,"method":"method1","params":[]}`,
			Notification{"method1", []interface{}{}},
			""},
		// 10
		{`{"id":null,"method":"method1","params":[1,"2"]}`,
			Notification{"method1", []interface{}{1, "2"}},
			""},
		// 11
		{`error`,
			nil,
			`invalid character 'e' looking for beginning of value`},
	} {
		buf := bytes.Buffer{}
		buf.WriteString(e.given)
		dec := NewDecoder(&buf)
		if e.err == "" {
			if d, err := dec.Decode(); err != nil {
				t.Error(err)
			} else if s := fmt.Sprintf("%#v", d); s != fmt.Sprintf("%#v", e.then) {
				t.Errorf("%v: %#v: %v", i, e.given, s)
			}
		} else {
			if _, err := dec.Decode(); err == nil {
				t.Error(err)
			} else if s := err.Error(); s != e.err {
				t.Error(s)
			}
		}
	}
}
