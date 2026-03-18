//go:build tamago && riscv64

package bn254_curve_dbl

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/bn254_curve_dbl.rs

type Point256 struct {
	X [4]uint64
	Y [4]uint64
}

//go:noinline
func Bn254Dbl(arr *Point256) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bn254_dbl(ptr)
}

func syscall_bn254_dbl(ptr uintptr)
