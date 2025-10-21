//go:build tamago && riscv64

package keccak

import (
	"unsafe"
)

//go:noinline
func Keccak(arr *[25]uint64) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_keccak(ptr)
}

func syscall_keccak(ptr uintptr)
