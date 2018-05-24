package buffers

import (
	"crypto/sha256"
	"hash/fnv"
)

// Buffers is a slice S of []byte that represent the single []byte that would
// result from concatenating them. N is the number of bytes in all []byte.
type Buffers struct {
	N int
	S [][]byte
}

// Append returns a new Buffers that represents the concatenation of all slices
// of f and those in q.
func (f Buffers) Append(q ...[]byte) Buffers {
	for _, p := range q {
		f = Buffers{f.N + len(p), append(f.S, p)}
	}
	return f
}

// Hash32 appends the FNV-1a hash of f to q and returns the resulting slice.
func (f Buffers) Hash32(q []byte) []byte {
	d := fnv.New32a()
	for _, p := range f.S {
		d.Write(p)
	}
	return d.Sum(q)
}

// Hash256 appends the SHA 256 hash of f to q and returns the resulting slice.
func (f Buffers) Hash256(q []byte) []byte {
	d := sha256.New()
	for _, p := range f.S {
		d.Write(p)
	}
	return d.Sum(q)
}
