package svb

// Uint32Decode256 vector decode stream vbytes with 256bits
func Uint32Decode256(masks, data []byte, out []uint32) {
	// bound check mask
	_ = masks[(len(out)+3)/4-1]
	i, offset := 0, 0
	for {
		cnt := len(masks) - i
		if cnt == 0 {
			break
		} else if cnt == 1 {
			// TODO: shuffle128
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
