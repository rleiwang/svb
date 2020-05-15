package perf

import (
	"testing"

	"svb"
)

func BenchmarkUint32Decode128(b *testing.B) {
	if !svb.HasXmm() {
		b.SkipNow()
	}

	for n := 0; n < b.N; n++ {
		svb.Uint32Decode128(masks, data, out)
	}
}

func BenchmarkUint32Decode256(b *testing.B) {
	if !svb.HasYmm() {
		b.SkipNow()
	}

	for n := 0; n < b.N; n++ {
		svb.Uint32Decode256(masks, data, out)
	}
}

func BenchmarkUint32Decode512(b *testing.B) {
	if !svb.HasZmm() {
		b.SkipNow()
	}

	for n := 0; n < b.N; n++ {
		svb.Uint32Decode512(masks, data, out)
	}
}
