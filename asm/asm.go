// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/buildtags"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func main() {
	Constraint(buildtags.Not("noasm").ToConstraint())
	Constraint(buildtags.Not("appengine").ToConstraint())
	Constraint(buildtags.Not("gccgo").ToConstraint())

	decode128()
	decode256()
	decode512()

	Generate()
}

type vectorShuffle interface {
	incr(i GP)
	cond(i GP, done LabelRef)
	body(m, i GP, lookup, offset Register, increment LabelRef)
}

type slice struct {
	ptr Register
	len Register
}

type registers struct {
	table Register
	masks slice
	data  slice
	out   slice
}

type xmmVector struct {
	registers
}

func (x *xmmVector) incr(i GP) {
	Comment("i++")
	LEAQ(Mem{Base: i.As64(), Disp: 1}, i.As64())
}

func (x *xmmVector) cond(i GP, done LabelRef) {
	Comment("i < len(masks)")
	CMPQ(i.As64(), x.masks.len)
	Comment("goto done if i >= len(masks)")
	JGE(done)
}

func (x *xmmVector) body(m, i GP, lookup, offset Register, increment LabelRef) {
	getMask(m, i, x.masks.ptr, 0)
	lookupShuffleMasks(lookup, x.table, m)

	step := GP64()
	MOVQ(i.As64(), step)
	Comment("step = i * 4 (4 integers)")
	SHLQ(U8(2), step)

	xmm := XMM()
	VMOVDQU(Mem{Base: x.data.ptr, Index: offset, Scale: 1}, xmm)
	PSHUFB(Mem{Base: lookup}, xmm)
	VMOVDQU(xmm, Mem{Base: x.out.ptr, Index: step, Scale: 4})

	incrementOffset(m, lookup, offset)
	JMP(increment)
}

type ymmVector struct {
	registers
}

func (y *ymmVector) incr(i GP) {
	Comment("i += 2")
	LEAQ(Mem{Base: i.As64(), Disp: 2}, i.As64())
}

func (y *ymmVector) cond(i GP, done LabelRef) {
	diff := GP64()
	MOVQ(y.masks.len, diff)
	SUBQ(i.As64(), diff)
	CMPQ(diff, U8(2))
	Comment("goto done if i >= len(masks)")
	JLT(done)
}

func (y *ymmVector) body(m, i GP, lookup, offset Register, increment LabelRef) {
	step := GP64()
	MOVQ(i.As64(), step)
	Comment("step = i * 4 (4 integers)")
	SHLQ(U8(2), step)

	maskYmm, dataYmm, outYmm := YMM(), YMM(), YMM()
	for k := 0; k < 2; k++ {
		Commentf("%dth DOUBLE QWORD", k)

		getMask(m, i, y.masks.ptr, k)
		lookupShuffleMasks(lookup, y.table, m)

		Commentf("move 16 bytes from ShuffleTable[masks[%d]] to %d double qword", k, k)
		VINSERTF128(U8(k), Mem{Base: lookup}, maskYmm, maskYmm)

		Commentf("move 16 bytes from data[offset] to %v double qword", k)
		VINSERTF128(U8(k), Mem{Base: y.data.ptr, Index: offset, Scale: 1}, dataYmm, dataYmm)

		incrementOffset(m, lookup, offset)
	}

	Comment("shuffle 8 uint32")
	VPSHUFB(maskYmm, dataYmm, outYmm)

	Comment("move 8 uint32 to out")
	VMOVDQU(outYmm, Mem{Base: y.out.ptr, Index: step, Scale: 4})

	JMP(increment)
}

type zmmVector struct {
	registers
}

func (z *zmmVector) incr(i GP) {
	Comment("i += 4")
	LEAQ(Mem{Base: i.As64(), Disp: 4}, i.As64())
}

func (z *zmmVector) cond(i GP, done LabelRef) {
	diff := GP64()
	MOVQ(z.masks.len, diff)
	SUBQ(i.As64(), diff)
	CMPQ(diff, U8(4))
	Comment("goto done if i >= len(masks)")
	JLT(done)
}

func (z *zmmVector) body(m, i GP, lookup, offset Register, increment LabelRef) {
	step := R8
	MOVQ(i.As64(), step)
	Comment("step = i * 4 (4 integers)")
	SHLQ(U8(2), step)

	expMask := R9W
	three := 3
	Commentf("init R9W expand mask %08b", three)
	MOVW(U16(three), expMask)

	for k := 0; k < 4; k++ {
		Commentf("%dth DOUBLE QWORD", k)
		if k > 0 {
			three <<= 2
			Commentf("expand mask R9 << 2, %08b", three)
			SHLW(U8(2), expMask)
		}

		getMask(m, i, z.masks.ptr, k)
		lookupShuffleMasks(lookup, z.table, m)

		Commentf(`AVX512, K1 = %08b
	KMOVW R9, K1`, three)

		Commentf(`AVX512, Move data[offset:] to Z0 with mask %08b
	VPEXPANDQ (SI)(R12*1), K1, Z0`, three)

		Commentf(`AVX512, Move ShuffleTable[masks[%d]] to Z1 with mask %08b  
	VPEXPANDQ (R11), K1, Z1`, k, three)

		incrementOffset(m, lookup, offset)
	}

	Comment(`AVX512, shuffle 16 uint32
	VPSHUFB Z1, Z0, Z2`)

	Comment(`AVX512, Copy 16 uint32 to out
	VMOVDQU8 Z2, (DI)(R8*4)`)

	JMP(increment)
}

func getMask(m, i GP, masks Register, d int) {
	Comment("m = masks[i]")
	MOVBQZX(Mem{Disp: d, Base: masks, Index: i.As64(), Scale: 1}, m.As64())
}

func getShuffleTable(shuffleTable Register) Register {
	Comment("shuffleTable = &ShuffleTable[256][16]")
	LEAQ(NewDataAddr(Symbol{Name: "·ShuffleTable"}, 0), shuffleTable)
	return shuffleTable
}

func lookupShuffleMasks(lookup, shuffleTable Register, m GP) {
	Comment("lookup = &ShuffleTable[m][16]")
	SHLQ(U8(4), m.As64())
	LEAQ(Mem{Base: shuffleTable, Index: m.As64(), Scale: 1}, lookup)
}

func incrementOffset(m GP, lookup, offset Register) {
	Comment("m >>= 6, note: m << 4 earlier")
	SHRL(U8(10), m.As32())
	Comment("m += 12")
	ADDL(U8(12), m.As32())
	Comment("lookup = ShuffleTable[m][12 + m >> 6]")
	MOVBQZX(Mem{Base: lookup, Index: m.As64(), Scale: 1}, lookup)
	Comment("offset += ShuffleTable[m][12 + m >> 6] + 1")
	LEAQ(Mem{Base: offset, Index: lookup, Scale: 1, Disp: 1}, offset)
}

func loop(suffix string, i, m GP, lookup, offset Register, vsf vectorShuffle) {
	increment := "increment" + suffix
	condition := "condition" + suffix
	done := "done" + suffix

	JMP(LabelRef(condition))

	// increment
	Label(increment)
	vsf.incr(i)

	// condition
	Label(condition)
	vsf.cond(i, LabelRef(done))

	// body
	vsf.body(m, i, lookup, offset, LabelRef(increment))

	// done
	Label(done)
}

func decode128() {
	TEXT("Uint32Decode128", NOSPLIT, "func(masks, data []byte, out []uint32)")
	Doc("Uint32Decode128 32 bits integer using XMM register, AVX")

	xmm := xmmVector{registers{
		table: getShuffleTable(GP64()),
		masks: slice{
			ptr: Load(Param("masks").Base(), GP64()),
			len: Load(Param("masks").Len(), GP64()),
		},
		data: slice{ptr: Load(Param("data").Base(), GP64())},
		out:  slice{ptr: Load(Param("out").Base(), GP64())},
	}}

	i := GP64()
	lookup := GP64()
	offset := GP64()
	m := GP64()

	XORQ(i, i)
	XORQ(m, m)
	XORQ(offset, offset)

	loop("", i, m, lookup, offset, &xmm)
	RET()
}

func decode256() {
	TEXT("Uint32Decode256", NOSPLIT, "func(masks, data []byte, out []uint32)")
	Doc("Uint32Decode256 32 bits integer using YMM register, AVX2")

	ymm := ymmVector{registers{
		table: getShuffleTable(GP64()),
		masks: slice{
			ptr: Load(Param("masks").Base(), GP64()),
			len: Load(Param("masks").Len(), GP64()),
		},
		data: slice{ptr: Load(Param("data").Base(), GP64())},
		out:  slice{ptr: Load(Param("out").Base(), GP64())},
	}}

	i := GP64()
	lookup := GP64()
	offset := GP64()
	m := GP64()

	XORQ(i, i)
	XORQ(m, m)
	XORQ(offset, offset)

	loop("_0", i, m, lookup, offset, &ymm)

	xmm := xmmVector{ymm.registers}
	loop("_1", i, m, lookup, offset, &xmm)

	RET()
}

func decode512() {
	TEXT("Uint32Decode512", NOSPLIT, "func(masks, data []byte, out []uint32)")
	Doc("Uint32Decode512 32 bits integer using ZMM register, AVX512")

	// use physical register since avo doesn't support AVX512 yet

	zmm := zmmVector{registers{
		table: getShuffleTable(RDX),
		masks: slice{
			ptr: Load(Param("masks").Base(), RAX),
			len: Load(Param("masks").Len(), RBX),
		},
		data: slice{ptr: Load(Param("data").Base(), RSI)},
		out:  slice{ptr: Load(Param("out").Base(), RDI)},
	}}

	// R8, R9
	i := R10
	lookup := R11
	offset := R12
	m := R13

	XORQ(i, i)
	XORQ(m, m)
	XORQ(offset, offset)

	loop("_0", i, m, lookup, offset, &zmm)

	xmm := xmmVector{zmm.registers}
	loop("_1", i, m, lookup, offset, &xmm)

	RET()
}
