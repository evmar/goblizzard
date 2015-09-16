package replay

import "bufio"

type bitReader struct {
	*bufio.Reader
	n   uint // number of bits remaining in buf
	buf uint8
}

func min(a, b uint) uint {
	if a < b {
		return a
	} else {
		return b
	}
}

func newBitReader(r *bufio.Reader) *bitReader {
	return &bitReader{Reader: r}
}

func (r *bitReader) ReadBits(n int) (uint64, error) {
	var val uint64
	for n > 0 {
		if r.n == 0 {
			b, err := r.ReadByte()
			if err != nil {
				return 0, err
			}
			r.buf = b
			r.n = 8
		}
		toCopy := min(r.n, uint(n))
		// This bit reading order is super confusing, see the test.
		mask := uint8((1 << toCopy) - 1)
		val = (val << toCopy) | uint64(r.buf&mask)
		r.n -= toCopy
		r.buf >>= toCopy
		n -= int(toCopy)
	}
	return val, nil
}

func (r *bitReader) SyncToByte() {
	r.n = 0
}
