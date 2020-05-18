package svb

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
