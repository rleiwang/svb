package svb

import (
	"math/rand"
	"testing"
)

func TestCodecGet(t *testing.T) {
	val := make([]uint32, 2838)
	for k := 0; k < len(val); k++ {
		val[k] = rand.Uint32() >> (31 & rand.Uint32())
	}

	codec := NewFromUint32(val)
	data := codec.Bytes()
	codec = NewFromBytes(data)
	t.Run("simple", func(t *testing.T) {
		for i, v := range val {
			if got := codec.Get(i); v != got {
				t.Fatalf("want %v, got %v\n", v, got)
			}
		}
	})
}
