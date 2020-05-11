package svb

import "testing"

var (
	data  = []byte{1, 2, 3, 4, 5, 6, 7, 8}
	masks = []byte{0, 0}
	out   = make([]uint32, 8)
)

func BenchmarkUint32Decode256(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Uint32Decode256(masks, data, out)
	}
}

func BenchmarkShuffle256(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Shuffle256(masks, data, 4, out)
	}
}
