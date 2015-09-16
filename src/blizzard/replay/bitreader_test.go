package replay

import (
	"bufio"
	"strings"
	"testing"
)

func expectBits(t *testing.T, exp uint64, r *bitReader, n int) {
	bits, err := r.ReadBits(n)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if bits != exp {
		t.Fatalf("expected %x, got %x", exp, bits)
	}
}

func TestBitReader(t *testing.T) {
	buf := bufio.NewReader(strings.NewReader("\x7c\x7c"))
	r := newBitReader(buf)
	// This is super confusing, but the bit reading happens from
	// the right side of each byte first.
	// 0111 1100 0111 1100
	//       ^^^
	expectBits(t, 0x4, r, 3)

	// 0111 1100 0111 1100
	//  ^^^ ^xxx
	expectBits(t, 0xf, r, 4)

	// 0111 1100 0111 1100
	// ^xxx xxxx       ^^^
	expectBits(t, 0x4, r, 4)

	// 0111 1100 0111 1100
	// xxxx xxxx ^^^^ ^xxx
	expectBits(t, 0xf, r, 5)
}
