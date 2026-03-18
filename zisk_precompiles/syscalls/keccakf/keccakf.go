//go:build tamago && riscv64

package keccakf

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/keccakf.rs

//go:noinline
func Keccak(arr *[25]uint64) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_keccakf(ptr)
}

func syscall_keccakf(ptr uintptr)
