package svb

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestUint32Decode32(t *testing.T) {
	val := make([]uint32, 35)
	for k := 0; k < len(val); k++ {
		val[k] = rand.Uint32() >> (31 & rand.Uint32())
	}

	masks, data = Uint32Encode(val)
	t.Run("simple", func(t *testing.T) {
		got := make([]uint32, len(val))
		Uint32Decode32(masks, data, got)
		if !reflect.DeepEqual(got, val) {
			t.Fatalf("want %v, got %v\n", val, got)
		}
	})
}
