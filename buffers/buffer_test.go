package buffers

import (
	"encoding/hex"
	"testing"
)

func TestBuffers_Hash32(t *testing.T) {
	h1 := Buffers{}.Append([]byte{0, 0}).Hash32(nil)
	if l := len(h1); l != 4 {
		t.Error(l)
	}
	if s := hex.EncodeToString(h1[:]); s != "117697cd" {
		t.Error(s)
	}
	h2 := Buffers{}.Append([]byte{0}, []byte{0}).Hash32(nil)
	if l := len(h2); l != 4 {
		t.Error(l)
	}
	if s := hex.EncodeToString(h2[:]); s != "117697cd" {
		t.Error(s)
	}
}

func TestBuffers_Hash256(t *testing.T) {
	h1 := Buffers{}.Append([]byte{0, 0}).Hash256(nil)
	if l := len(h1); l != 32 {
		t.Error(l)
	}
	if s := hex.EncodeToString(h1[:]); s != "96a296d224f285c67bee93c30f8a309157f0daa35dc5b87e410b78630a09cfc7" {
		t.Error(s)
	}
	h2 := Buffers{}.Append([]byte{0}, []byte{0}).Hash256(nil)
	if l := len(h2); l != 32 {
		t.Error(l)
	}
	if s := hex.EncodeToString(h2[:]); s != "96a296d224f285c67bee93c30f8a309157f0daa35dc5b87e410b78630a09cfc7" {
		t.Error(s)
	}
}
