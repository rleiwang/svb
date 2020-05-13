package svb

import (
	"reflect"
	"testing"
)

func TestEncodeUint32(t *testing.T) {
	type args struct {
		value []uint32
	}
	tests := []struct {
		name     string
		args     args
		wantMask []byte
		wantData []byte
	}{
		{"basic", args{[]uint32{1 << 28, 1 << 20, 1 << 12, 1 << 4}}, []byte{0b00011011}, []byte{0, 0, 0, 16, 0, 0, 16, 0, 16, 16}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMask, gotData := Uint32Encode(tt.args.value)
			if !reflect.DeepEqual(gotMask, tt.wantMask) {
				t.Errorf("EncodeUint32() gotMask = %v, want %v", gotMask, tt.wantMask)
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("EncodeUint32() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}
