//go:build tamago && riscv64

package secp256r1_dbl

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/secp256r1_dbl.rs

type Point256 struct {
	X [4]uint64
	Y [4]uint64
}

//go:noinline
func Secp256r1Dbl(arr *Point256) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_secp256r1_dbl(ptr)
}

func syscall_secp256r1_dbl(ptr uintptr)
