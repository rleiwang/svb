package svb

import "github.com/intel-go/cpuid"

func HasZmm() bool {
	return cpuid.EnabledAVX512 && cpuid.HasExtendedFeature(cpuid.AVX512F)
}

func HasYmm() bool {
	return cpuid.EnabledAVX && cpuid.HasExtendedFeature(cpuid.AVX2)
}

func HasXmm() bool {
	return cpuid.HasFeature(cpuid.SSSE3) && cpuid.HasFeature(cpuid.AVX)
}
