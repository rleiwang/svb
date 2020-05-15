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
	for ; len(masks)-i >= 2; i += 2 {
		len := ShuffleTable[masks[i]][12+masks[i]>>6] + 1
		len += ShuffleTable[masks[i+1]][12+masks[i+1]>>6] + 1
		offset += int(Shuffle256(masks[i:], data[offset:], out[i*4:]))
	}
	if len(masks) > i {
		Shuffle128(ShuffleTable[masks[i]][:], data[offset:], out[i*4:])
	}
}

// Uint32Decode512 vector decode stream vbytes with 512 bits vector registers
func Uint32Decode512(masks, data []byte, out []uint32) {
	// bound check mask
	_ = masks[(len(out)+3)/4-1]
	i, offset := 0, 0
	for ; len(masks)-i >= 4; i += 4 {
		offset += int(Shuffle512(masks[i:], data[offset:], out[i*4:]))
	}
	if len(masks) > i {
		Uint32Decode256(masks[i:], data[offset:], out[i*4:])
	}
}
