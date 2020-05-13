package svb

import "encoding/binary"

// Uint32Encode encode []uint32 with stream vbytes codec as LittleEndian
func Uint32Encode(value []uint32) (mask, data []byte) {
	mask = make([]byte, (len(value)+3)/4)
	data = make([]byte, len(value)*4)
	cnt := 0

	for i, v := range value {
		code, buf := byte(0), make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, v)
		if v < 1<<8 {
			copy(data[cnt:], buf[:1])
		} else if v < 1<<16 {
			copy(data[cnt:], buf[:2])
			code = 1
		} else if v < 1<<24 {
			copy(data[cnt:], buf[:3])
			code = 2
		} else {
			copy(data[cnt:], buf)
			code = 3
		}
		// LittleEndian
		mask[i/4] |= code << ((i % 4) << 1)
		cnt += int(code) + 1
	}

	data = data[:cnt]
	return
}
