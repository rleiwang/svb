package svb

// Uint32Decode128 vector decode stream vbytes with 128 bits vector registers
func Uint32Decode128(masks, data []byte, out []uint32) {
	// bound check mask
	_ = masks[(len(out)+3)/4-1]
	offset := 0
	for i := 0; i < len(masks); i++ {
		m := masks[i]
		Shuffle128(ShuffleTable[m][:], data[offset:], out[i*4:])
		offset += int(ShuffleTable[m][12+m>>6]) + 1
	}
}

// Uint32Decode256 vector decode stream vbytes with 256 bits vector registers
func Uint32Decode256(masks, data []byte, out []uint32) {
	// bound check mask
	_ = masks[(len(out)+3)/4-1]
	i, offset := 0, 0
	for {
		cnt := len(masks) - i
		if cnt == 0 {
			break
		} else if cnt == 1 {
			Shuffle128(ShuffleTable[masks[i]][:], data[offset:], out[i*4:])
			break
		}
		upper := masks[i]
		upperOffset := int(ShuffleTable[upper][12+upper>>6]) + 1

		Shuffle256(masks[i:], data[offset:], upperOffset, out[i*4:])

		offset += upperOffset
		lower := masks[i+1]
		offset += int(ShuffleTable[lower][12+lower>>6]) + 1
		i += 2
	}
}
