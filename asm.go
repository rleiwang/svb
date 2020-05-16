// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func main() {
	decode128()
	shuffle128()
	shuffle256()
	shuffle512()

	Generate()
}

func shuffle128() {
	TEXT("Shuffle128", NOSPLIT, "func(shuffle, data []byte, out []uint32)")
	Doc("Shuffle 32 bits integer using XMM register, AVX")
	shufflePtr := Load(Param("shuffle").Base(), GP64())
	dataPtr := Load(Param("data").Base(), GP64())
	outPtr := Load(Param("out").Base(), GP64())
	xmm := XMM()
	VMOVDQU(Mem{Base: dataPtr}, xmm)
	PSHUFB(Mem{Base: shufflePtr}, xmm)
	VMOVDQU(xmm, Mem{Base: outPtr})
	RET()
}

func decode128() {
	TEXT("Uint32Decode128", NOSPLIT, "func(masks, data []byte, out []uint32)")
	Doc("Uint32Decode128 32 bits integer using XMM register, AVX")
	masksPtr := Load(Param("masks").Base(), GP64())
	masksLen := Load(Param("masks").Len(), GP64())
	dataPtr := Load(Param("data").Base(), GP64())
	outPtr := Load(Param("out").Base(), GP64())

	condition := "condition"
	increment := "increment"
	done := "done"

	Label("init")
	Comment("i := 0")
	i := GP64()
	XORQ(i, i)

	Comment("offset := 0")
	offset := GP64()
	XORQ(offset, offset)

	Comment("shuffleTable = &ShuffleTable[256][16]")
	shuffleTable := GP64()
	LEAQ(NewDataAddr(Symbol{Name: "·ShuffleTable"}, 0), shuffleTable)

	Comment("var lookup = &ShuffleTable[masks[i]]")
	lookup := GP64()
	JMP(LabelRef(condition))

	// increment
	Label(increment)
	Comment("i++")
	LEAQ(Mem{Base: i, Disp: 1}, i)

	Label(condition)
	Comment("i < len(masks)")
	CMPQ(i, masksLen)
	Comment("goto done if i >= len(masks)")
	JGE(LabelRef(done))

	Comment("body")
	m := GP64()
	Comment("m = masks[i]")
	MOVBQZX(Mem{Base: masksPtr, Index: i, Scale: 1}, m)

	Comment("lookup = &ShuffleTable[m][16]")
	SHLQ(U8(4), m)
	LEAQ(Mem{Base: shuffleTable, Index: m, Scale: 1}, lookup)

	step := GP64()
	MOVQ(i, step)
	Comment("step = i * 4 (4 integers)")
	SHLQ(U8(4), step)

	xmm := XMM()
	VMOVDQU(Mem{Base: dataPtr, Index: offset, Scale: 1}, xmm)
	PSHUFB(Mem{Base: lookup}, xmm)
	VMOVDQU(xmm, Mem{Base: outPtr, Index: step, Scale: 1})

	Comment("m >>= 6, note: m << 4 earlier")
	SHRL(U8(10), m.As32())
	Comment("m += 12")
	ADDW(U8(12), m.As16())
	Comment("lookup = ShuffleTable[m][12 + m >> 6]")
	MOVBQZX(Mem{Base: lookup, Index: m, Scale: 1}, lookup)
	Comment("offset += ShuffleTable[m][12 + m >> 6] + 1")
	LEAQ(Mem{Base: offset, Index: lookup, Scale: 1, Disp: 1}, offset)
	JMP(LabelRef(increment))

	Label(done)
	RET()
}

func shuffle256() {
	TEXT("Shuffle256", NOSPLIT, "func(masks, data []byte, out []uint32) byte")
	Doc("Shuffle 32 bits integer using YMM register, AVX2")
	masksPtr := Load(Param("masks").Base(), GP64())
	dataPtr := Load(Param("data").Base(), GP64())
	outPtr := Load(Param("out").Base(), GP64())

	Comment("&ShuffleTable[256][16]")
	shuffleTable := Mem{Base: GP64()}
	LEAQ(NewDataAddr(Symbol{Name: "·ShuffleTable"}, 0), shuffleTable.Base)

	offset := GP64()
	Comment("offset := 0")
	XORQ(offset, offset)

	maskYmm, dataYmm, outYmm := YMM(), YMM(), YMM()
	for y := 0; y < 2; y++ {
		Commentf("%dth DOUBLE QWORD", y)

		maskOffset := GP64()
		Commentf("r = masks[%d] ", y)
		MOVBQZX(Mem{Disp: y, Base: masksPtr, Scale: 1}, maskOffset)

		Comment("shuffle table is [256][16], offset *= 16, left shift 4 bits")
		SHLQ(U8(4), maskOffset)

		st := GP64()
		Commentf("R = &ShuffleTable[masks[%d]]", y)
		LEAQ(Mem{Base: shuffleTable.Base, Index: maskOffset, Scale: 1}, st)

		Commentf("move 16 bytes from ShuffleTable[masks[%d]] to %d double qword", y, y)
		VINSERTF128(U8(y), Mem{Base: st}, maskYmm, maskYmm)

		Commentf("move 16 bytes from data[offset] to %v double qword", y)
		VINSERTF128(U8(y), Mem{Base: dataPtr, Index: offset, Scale: 1}, dataYmm, dataYmm)

		Comment("maskOffset >> 10, as m >> 6")
		SHRQ(U8(10), maskOffset)
		Comment("m += 12")
		LEAQ(Mem{Base: maskOffset, Disp: 12}, maskOffset)
		Comment("v = ShuffleTable[key][12 + key >> 6]")
		MOVBQZX(Mem{Base: st, Index: maskOffset, Scale: 1}, st)

		Comment("data offset += v + 1")
		LEAQ(Mem{Base: st, Index: offset, Scale: 1, Disp: 1}, offset)
	}

	Comment("shuffle 8 uint32")
	VPSHUFB(maskYmm, dataYmm, outYmm)

	Comment("move 8 uint32 to out")
	VMOVDQU(outYmm, Mem{Base: outPtr})

	Store(offset.As8(), ReturnIndex(0))
	RET()
}

func shuffle512() {
	TEXT("Shuffle512", NOSPLIT, "func(masks, data []byte, out []uint32) byte")
	Doc("Shuffle 32 bits integer using ZMM register, AVX512")
	// use physical register since avo doesn't support AVX512 yet
	masksPtr := Load(Param("masks").Base(), RAX)
	Load(Param("data").Base(), RBX)
	Load(Param("out").Base(), RCX)

	//_, _, _ = ZMM(), ZMM(), ZMM()
	offset := R8
	Comment("Clear data offset, R8")
	XORQ(offset, offset)

	expMask := R9W
	three := 3
	Commentf("init R9W expand mask %08b", three)
	MOVW(U16(three), expMask)

	shuffleTable := RDX
	Comment("DX = &ShuffleTable[256][16]")
	LEAQ(NewDataAddr(Symbol{Name: "·ShuffleTable"}, 0), Mem{Base: shuffleTable}.Base)

	si, st := RSI, R10
	for i := 0; i < 4; i++ {
		Commentf("%dth DOUBLE QWORD", i)
		if i > 0 {
			three <<= 2
			Commentf("expand mask R9 << 2, %08b", three)
			SHLW(U8(2), expMask)
		}

		Commentf(`AVX512, K1 = %08b
	KMOVW R9, K1`, three)

		Commentf(`AVX512, Move data[offset:] to Z0 with mask %08b
	VPEXPANDQ (BX)(R8*1), K1, Z0`, three)

		Commentf("SI = masks[%d]", i)
		MOVBQZX(Mem{Base: masksPtr, Disp: i}, si)
		Comment("<< 4 bits, 16 bytes offset, ShuffleTable[256][16]")
		SHLQ(U8(4), si)
		Commentf("R10 = ShuffleTable[mask[%d]]", i)
		LEAQ(Mem{Base: shuffleTable, Index: si, Scale: 1}, st)

		Commentf(`AVX512, Move ShuffleTable[masks[%d]] to Z1 with mask %08b  
	VPEXPANDQ (R10), K1, Z1`, i, three)

		Comment("SI >> 10, as m >> 6")
		SHRQ(U8(10), si)
		tmp := R11
		Commentf("R11 = ShuffleTable[SI][12+SI>>6], SI = masks[%d]", i)
		MOVBQZX(Mem{Base: st, Index: si, Disp: 12, Scale: 1}, tmp)

		Comment("data offset += R11 + 1")
		LEAQ(Mem{Base: tmp, Index: offset, Scale: 1, Disp: 1}, offset)
	}

	Comment(`AVX512, shuffle 16 uint32
	VPSHUFB Z1, Z0, Z2`)

	Comment(`AVX512, Copy 16 uint32 to out
	VMOVDQU8 Z2, (CX)`)

	Comment("Return processed data length")
	Store(offset.As8(), ReturnIndex(0))
	RET()
}
