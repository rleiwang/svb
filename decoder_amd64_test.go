package svb

import (
	"math/rand"
	"reflect"
	"testing"
)

var (
	data  = []byte{1, 2, 3, 4, 5, 6, 7, 8}
	masks = []byte{0, 0}
)

func TestUint32Decode128(t *testing.T) {
	if !HasXmm() {
		t.SkipNow()
	}

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
	if !HasYmm() {
		t.SkipNow()
	}

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

func TestUint32Decode512(t *testing.T) {
	if !HasZmm() {
		t.SkipNow()
	}

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
		{"simple",
			args{
				masks: []byte{0, 0, 0, 0},
				data:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				n:     16,
			},
			[]uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]uint32, tt.args.n)
			Uint32Decode512(tt.args.masks, tt.args.data, got)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}

func TestEncodeDecode128Uint32(t *testing.T) {
	if !HasXmm() {
		t.SkipNow()
	}

	t.Run("encode_decode", func(t *testing.T) {
		want := make([]uint32, 48)
		for k := 0; k < len(want); k++ {
			want[k] = rand.Uint32() >> (31 & rand.Uint32())
		}

		masks, data = Uint32Encode(want)

		got := make([]uint32, len(want))
		Uint32Decode128(masks, data, got)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Encode = %v, Decode = %v", want, got)
		}
	})
}

func TestEncodeDecode256Uint32(t *testing.T) {
	if !HasYmm() {
		t.SkipNow()
	}

	t.Run("encode_decode", func(t *testing.T) {
		want := make([]uint32, 56)
		for k := 0; k < len(want); k++ {
			want[k] = rand.Uint32() >> (31 & rand.Uint32())
		}

		masks, data = Uint32Encode(want)

		got := make([]uint32, len(want))
		Uint32Decode256(masks, data, got)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Encode = %v, Decode = %v", want, got)
		}
	})
}

func TestEncodeDecode512Uint32(t *testing.T) {
	if !HasZmm() {
		t.SkipNow()
	}

	t.Run("encode_decode", func(t *testing.T) {
		want := make([]uint32, 63)
		for k := 0; k < len(want); k++ {
			want[k] = rand.Uint32() >> (31 & rand.Uint32())
		}

		masks, data = Uint32Encode(want)

		got := make([]uint32, len(want))
		Uint32Decode512(masks, data, got)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Encode = %v, Decode = %v", want, got)
		}
	})
}
