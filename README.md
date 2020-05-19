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
|Uint32Decode32|13562586ns|17394533ns|
|Uint32Decode128|331368ns|411102ns|
|Uint32Decode256|327100ns|406230ns|
|Uint32Decode512|470571ns|569497ns|
