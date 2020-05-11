package svb

import (
	"reflect"
	"testing"
)

var (
	data  = []byte{1, 2, 3, 4, 5, 6, 7, 8}
	masks = []byte{0, 0}
)

func TestUint32Decode128(t *testing.T) {
	type args struct {
		masks []byte
		data  []byte
		n     int
	}
	tests := []struct {
		name string
		args args
		want []uint32
	}{
		{"simple", args{masks: masks, data: data, n: 8}, []uint32{1, 2, 3, 4, 5, 6, 7, 8}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]uint32, tt.args.n)
			Uint32Decode128(tt.args.masks, tt.args.data, got)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}

func TestUint32Decode256(t *testing.T) {
	type args struct {
		masks []byte
		data  []byte
		n     int
	}
	tests := []struct {
		name string
		args args
		want []uint32
	}{
		{"simple", args{masks: masks, data: data, n: 8}, []uint32{1, 2, 3, 4, 5, 6, 7, 8}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]uint32, tt.args.n)
			Uint32Decode256(tt.args.masks, tt.args.data, got)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}

func BenchmarkUint32Decode128(b *testing.B) {
	out := make([]uint32, 8)
	for n := 0; n < b.N; n++ {
		Uint32Decode128(masks, data, out)
	}
}

func BenchmarkShuffle128(b *testing.B) {
	out := make([]uint32, 8)
	for n := 0; n < b.N; n++ {
		Shuffle128(ShuffleTable[masks[0]][:], data, out)
		Shuffle128(ShuffleTable[masks[1]][:], data, out[4:])
	}
}

func BenchmarkUint32Decode256(b *testing.B) {
	out := make([]uint32, 8)
	for n := 0; n < b.N; n++ {
		Uint32Decode256(masks, data, out)
	}
}

func BenchmarkShuffle256(b *testing.B) {
	out := make([]uint32, 8)
	for n := 0; n < b.N; n++ {
		Shuffle256(masks, data, 4, out)
	}
}
