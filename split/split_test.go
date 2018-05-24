package split

import (
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/daniel-fanjul-alcuten/floc/buffers"
)

var testSplit = Split{0, 1 << 20, 1<<12 - 1, 1<<12 - 1, 1 << 13, nil}

func TestSplit_Nil(t *testing.T) {
	s := testSplit
	g, r := s.Split(buffers.Buffers{})
	if n := g.N; n != 0 {
		t.Error(n)
	} else if n := r.N; n != 0 {
		t.Error(n)
	}
}

func TestSplit_Empty(t *testing.T) {
	s := testSplit
	g, r := s.Split(buffers.Buffers{}.Append([]byte{}))
	if n := g.N; n != 0 {
		t.Error(n)
	} else if n := r.N; n != 0 {
		t.Error(n)
	}
}

func TestSplit_Mask(t *testing.T) {
	s := testSplit
	s.Mask, s.Cond = 1<<3-1, 1<<3-1
	g, r := s.Split(buffers.Buffers{}.Append([]byte("\xff!")))
	if n := g.N; n != 1 {
		t.Error(n)
	} else if n := r.N; n != 1 {
		t.Error(n)
	}
	g, r = s.Split(buffers.Buffers{}.Append([]byte("\xf0!")))
	if n := g.N; n != 2 {
		t.Error(n)
	} else if n := r.N; n != 0 {
		t.Error(n)
	}
	g, r = s.Split(buffers.Buffers{}.Append([]byte("\xf0\x0f!")))
	if n := g.N; n != 2 {
		t.Error(n)
	} else if n := r.N; n != 1 {
		t.Error(n)
	}
	g, r = s.Split(buffers.Buffers{}.Append([]byte("\xf0\x00!")))
	if n := g.N; n != 3 {
		t.Error(n)
	} else if n := r.N; n != 0 {
		t.Error(n)
	}
}

func TestSplit_Min(t *testing.T) {
	s := testSplit
	s.Mask, s.Cond = 1<<3-1, 1<<3-1
	g, r := s.Split(buffers.Buffers{}.Append([]byte("\xff!")))
	if n := g.N; n != 1 {
		t.Error(n)
	} else if n := r.N; n != 1 {
		t.Error(n)
	}
	s.Min = 2
	g, r = s.Split(buffers.Buffers{}.Append([]byte("\xff!")))
	if n := g.N; n != 2 {
		t.Error(n)
	} else if n := r.N; n != 0 {
		t.Error(n)
	}
}

func TestSplit_Max(t *testing.T) {
	s := testSplit
	g, r := s.Split(buffers.Buffers{}.Append([]byte("abcde")))
	if n := g.N; n != 5 {
		t.Error(n)
	} else if n := r.N; n != 0 {
		t.Error(n)
	}
	s.Max = 3
	g, r = s.Split(buffers.Buffers{}.Append([]byte("abcde")))
	if n := g.N; n != 3 {
		t.Error(n)
	} else if n := r.N; n != 2 {
		t.Error(n)
	}
}

var benchmarkSplitBuffer []byte

func BenchmarkSplit(t *testing.B) {
	if benchmarkSplitBuffer == nil {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		p := make([]byte, 1024*1024*1024)
		if _, err := io.ReadFull(r, p); err != nil {
			t.Fatal(err)
		}
		benchmarkSplitBuffer = p
	}
	s := Split{0, 1 << 20, 1<<12 - 1, 1<<12 - 1, 1 << 13, nil}
	s.Reset()
	t.ResetTimer()
	for n := t.N; n > 0; {
		p := benchmarkSplitBuffer
		if len(p) > n {
			p = p[:n]
		}
		n -= len(p)
		for f := (buffers.Buffers{}).Append(p); f.N > 0; {
			_, f = s.Split(f)
		}
	}
	t.SetBytes(int64(t.N))
}
