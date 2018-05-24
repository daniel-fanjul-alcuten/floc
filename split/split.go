package split

import "github.com/daniel-fanjul-alcuten/floc/buffers"

// Split is the struct with the configuration and scratch space for the the
// rolling hash implemented in the Split function.
type Split struct {

	// Min is the mininum allowed length of the Split.
	Min int

	// Max is the maximum allowed length of the Split.
	Max int

	// Mask is the mask applied to the current value of the rolling hash.
	Mask uint32

	// Cond is the value that the current value of the rolling hash must have
	// after applying the Mask for the Split to trigger.
	Cond uint32

	// Window is the size of the slice Ring.
	Window int

	// Ring is the scratch space reserved for the window of bytes of the rolling
	// hash.
	Ring []uint8
}

// Reset sets the configuration fields that are zero values to proper default
// values when these zero values are not valid values. It allocates or zeros
// the Ring.
func (h *Split) Reset() {
	if h.Max == 0 {
		h.Max = 1<<31 - 1
	}
	if h.Mask == 0 {
		h.Mask = 1<<12 - 1
	}
	if h.Cond == 0 {
		h.Cond = 1<<12 - 1
	}
	if h.Window == 0 {
		h.Window = 1 << 13
	}
	if h.Ring == nil || len(h.Ring) != h.Window {
		h.Ring = make([]uint8, h.Window)
	} else {
		for j := range h.Ring {
			h.Ring[j] = 0
		}
	}
}

// Split splits the Buffers f into the Buffers g and r such that the
// concatenation of g and r is equal to f, g.N >= h.Min, g.N <= h.Max, and the
// rolling hash of g masked by h.Mask is equal to h.Cond or else g.N == h.Max.
// h.Reset() is invoked before anything else.
func (h *Split) Split(f buffers.Buffers) (g buffers.Buffers, r buffers.Buffers) {
	h.Reset()
	l, acc, min, max := 0, uint32(0), h.Min, h.Max
	for i, p := range f.S {
		for j := 0; j < len(p); {
			for ; j+3 < len(p) && l+3 < len(h.Ring); j += 4 {
				acc += uint32(p[j+0]) - uint32(h.Ring[l+0])
				if (j+1 >= min && acc&h.Mask == h.Cond) || j+1 >= max {
					g = g.Append(p[:j+1])
					r = r.Append(p[j+1:]).Append(f.S[i+1:]...)
					return
				}
				acc += uint32(p[j+1]) - uint32(h.Ring[l+1])
				if (j+2 >= min && acc&h.Mask == h.Cond) || j+2 >= max {
					g = g.Append(p[:j+2])
					r = r.Append(p[j+2:]).Append(f.S[i+1:]...)
					return
				}
				acc += uint32(p[j+2]) - uint32(h.Ring[l+2])
				if (j+3 >= min && acc&h.Mask == h.Cond) || j+3 >= max {
					g = g.Append(p[:j+3])
					r = r.Append(p[j+3:]).Append(f.S[i+1:]...)
					return
				}
				acc += uint32(p[j+3]) - uint32(h.Ring[l+3])
				if (j+4 >= min && acc&h.Mask == h.Cond) || j+4 >= max {
					g = g.Append(p[:j+4])
					r = r.Append(p[j+4:]).Append(f.S[i+1:]...)
					return
				}
				h.Ring[l+0] = p[j+0]
				h.Ring[l+1] = p[j+1]
				h.Ring[l+2] = p[j+2]
				h.Ring[l+3] = p[j+3]
				l = (l + 4) % len(h.Ring)
			}
			if j < len(p) {
				acc += uint32(p[j]) - uint32(h.Ring[l])
				if (j+1 >= min && acc&h.Mask == h.Cond) || j+1 >= max {
					g = g.Append(p[:j+1])
					r = r.Append(p[j+1:]).Append(f.S[i+1:]...)
					return
				}
				h.Ring[l], l = p[j], (l+1)%len(h.Ring)
				j++
			}
		}
		g = g.Append(p)
		min, max = min-len(p), max-len(p)
	}
	return
}
