// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	// Shuffle128
	TEXT("Shuffle128", NOSPLIT, "func(shuffle, data []byte, out []uint32)")
	Doc("Shuffle 32 bits integer with XMM register")
	shufflePtr := Load(Param("shuffle").Base(), GP64())
	dataPtr := Load(Param("data").Base(), GP64())
	outPtr := Load(Param("out").Base(), GP64())
	xmm := XMM()
	VMOVDQU(Mem{Base: dataPtr}, xmm)
	PSHUFB(Mem{Base: shufflePtr}, xmm)
	VMOVDQU(xmm, Mem{Base: outPtr})
	RET()

	// Shuffle256
	TEXT("Shuffle256", NOSPLIT, "func(masks, data []byte, offset int, out []uint32)")
	Doc("Shuffle 32 bits integer with YMM register")
	masksPtr := Load(Param("masks").Base(), GP64())
	dataPtr = Load(Param("data").Base(), GP64())
	outPtr = Load(Param("out").Base(), GP64())

	Comment("ShuffleTable[256][16]")
	shuffleTable := Mem{Base: GP64()}
	LEAQ(NewDataAddr(Symbol{Name: "·ShuffleTable"}, 0), shuffleTable.Base)

	offset := GP64()
	XORQ(offset, offset)

	words := []string{"lower", "upper"}
	maskYmm, dataYmm, outYmm := YMM(), YMM(), YMM()
	for y := 0; y < 2; y++ {
		maskOffset := GP64()
		Commentf("move mask %v byte to GPR", words[y])
		MOVBQZX(Mem{Disp: y, Base: masksPtr, Scale: 1}, maskOffset)

		Comment("shuffle table is [256][16], offset *= 16, left shift 4 bits")
		SHLQ(U8(4), maskOffset)
		Commentf("move 16 shuffle bytes from ShuffleTable to %v dword", words[y])
		VINSERTF128(U8(y), Mem{Base: shuffleTable.Base, Index: maskOffset, Scale: 1}, maskYmm, maskYmm)

		Commentf("move 128 bits (16 bytes) data bytes to %v dword", words[y])
		VINSERTF128(U8(y), Mem{Base: dataPtr, Index: offset, Scale: 1}, dataYmm, dataYmm)

		Comment("load data offset")
		Load(Param("offset"), offset)
	}

	Comment("shuffle 8 uint32")
	VPSHUFB(maskYmm, dataYmm, outYmm)

	Comment("move 8 uint32 to out")
	VMOVDQU(outYmm, Mem{Base: outPtr})
	RET()

	Generate()
}
