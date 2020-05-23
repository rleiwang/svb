package svb

import (
	"encoding/binary"
)

const (
	// blockSZ must be blockSZ % 4 == 0
	blockSZ = 512
)

type Codec struct {
	mask []byte
	data []byte
	part []int
	cnt  int
}

func NewFromUint32(data []uint32) *Codec {
	codec := &Codec{
		cnt: len(data),
	}
	codec.mask, codec.data = Uint32Encode(data)
	codec.part = partition(codec.mask, len(data))
	return codec
}

func NewFromBytes(data []byte) *Codec {
	codec := &Codec{
		cnt: int(binary.LittleEndian.Uint32(data[:4])),
	}
	offset := (codec.cnt+3)/4 + 4
	codec.mask = data[4:offset]
	codec.data = data[offset:]
	codec.part = partition(codec.mask, codec.cnt)

	return codec
}

func (c *Codec) Len() int {
	return c.cnt
}

func (c *Codec) Bytes() []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(c.cnt))
	return append(buf[:], append(c.mask, c.data...)...)
}

// Get 0-based ith uint32
func (c *Codec) Get(i int) uint32 {
	nth, offset := i/blockSZ, 0
	if nth > 0 {
		offset = c.part[nth-1]
	}

	nth *= blockSZ / 4
	last := i / 4
	for ; nth < last; nth++ {
		m := c.mask[nth]
		offset += int(ShuffleTable[m][12+m>>6]) + 1
	}
	m := c.mask[last]
	for last *= 4; last < i; last++ {
		offset += int(m&byte(3)) + 1
		m >>= 2
	}
	step := int(m&byte(3)) + 1
	buf := []byte{0, 0, 0, 0}
	copy(buf, c.data[offset:(offset+step)])

	return binary.LittleEndian.Uint32(buf)
}

func partition(masks []byte, cnt int) []int {
	parts := cnt / blockSZ
	if parts == 0 {
		return nil
	}

	partitions := make([]int, parts)
	offset, lastp := 0, 0
	parts = blockSZ / 4
	for i := 0; i < len(masks); i++ {
		p := i / parts
		if p > lastp {
			partitions[lastp] = offset
			lastp = p
		}
		m := masks[i]
		offset += int(ShuffleTable[m][12+m>>6]) + 1
	}
	return partitions
}
