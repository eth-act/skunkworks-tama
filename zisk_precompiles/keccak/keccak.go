//go:build tamago && riscv64

package keccak

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/blob/0207f975d0f03724f6f5784692ee2725e0d340d3/ziskos/entrypoint/src/syscalls/keccakf.rs

//go:noinline
func Keccak(arr *[25]uint64) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_keccak(ptr)
}

func syscall_keccak(ptr uintptr)
