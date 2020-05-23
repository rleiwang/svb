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

	codec *svb.Codec

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

	codec = svb.NewFromUint32(val)

	os.Exit(m.Run())
}

func BenchmarkUint32Decode32(b *testing.B) {
	for n := 0; n < b.N; n++ {
		svb.Uint32Decode32(masks, data, out)
	}
}

func BenchmarkCodecGet(b *testing.B) {
	pos := make([]int, 1000)
	for i := range pos {
		pos[i] = int(rand.Uint32()) % *sz
	}

	svb.Uint32Decode(masks, data, out)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, p := range pos {
			if codec.Get(p) != out[p] {
				b.Fatalf("want %v, got %v @%d", out[p], codec.Get(p), p)
			}
		}
	}
}
