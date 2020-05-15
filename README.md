# Stream VByte in Go with SIMD
This is another pure Go implementation of [Stream VByte: Faster Byte-Oriented Integer Compression](https://arxiv.org/abs/1709.08990).

It uses [avo](https://github.com/mmcloughlin/avo) by Michael McLoughlin to generate Go assembler code. This Go implementation has referenced https://github.com/lemire/streamvbyte

### Speed Test

performance benchmark measures the latency to decode 1 millions uint32.

```bash
svb/perf ❯❯❯ go test -bench .
```

|Function|Cascade Lake|Skylake|
|---|---|---|
|Uint32Decode32|12358406ns|17420943ns|
|Uint32Decode128|1423395ns|1988929ns|
|Uint32Decode256|1111386ns|1533913ns|
|Uint32Decode512|810798ns|1095372ns|
