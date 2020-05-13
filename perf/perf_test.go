package perf

import (
	"flag"
	"math/rand"
	"os"
	"testing"

	"svb"
)

var (
	masks []byte
	data  []byte
	out   []uint32

	sz = flag.Int("sz", 1000000, "number of uint32, default is 1M")
)

func TestMain(m *testing.M) {
	flag.Parse()

	val := make([]uint32, *sz)
	for k := 0; k < len(val); k++ {
		val[k] = rand.Uint32() >> (31 & rand.Uint32())
	}

	masks, data = svb.Uint32Encode(val)
	out = make([]uint32, len(val))

	os.Exit(m.Run())
}

func BenchmarkUint32Decode128(b *testing.B) {
	for n := 0; n < b.N; n++ {
		svb.Uint32Decode128(masks, data, out)
	}
}

func BenchmarkUint32Decode256(b *testing.B) {
	for n := 0; n < b.N; n++ {
		svb.Uint32Decode256(masks, data, out)
	}
}

func BenchmarkUint32Decode512(b *testing.B) {
	for n := 0; n < b.N; n++ {
		svb.Uint32Decode512(masks, data, out)
	}
}
