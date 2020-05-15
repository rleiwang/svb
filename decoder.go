package svb

import "encoding/binary"

var (
	decoder func([]byte, []byte, []uint32)
)

func init() {
	if HasZmm() {
		decoder = Uint32Decode512
	} else if HasYmm() {
		decoder = Uint32Decode256
	} else if HasXmm() {
		decoder = Uint32Decode128
	} else {
		decoder = Uint32Decode32
	}
}

func Uint32Decode(masks, data []byte, out []uint32) {
	decoder(masks, data, out)
}

func Uint32Decode32(masks, data []byte, out []uint32) {
	offset := 0
	buf := make([]byte, 4)
	for i, j := 0, 0; i < len(masks); i++ {
		mask := byte(3)
		for k := 0; k < 4 && j < len(out); k++ {
			step := int((masks[i]&mask)>>(k*2)) + 1
			copy(buf, []byte{0, 0, 0, 0})
			copy(buf, data[offset:(offset+step)])
			out[j] = binary.LittleEndian.Uint32(buf)
			j++
			offset += step
			mask <<= 2
		}
	}
}
