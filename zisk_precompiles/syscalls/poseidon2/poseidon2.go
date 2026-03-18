//go:build tamago && riscv64

package poseidon2

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/poseidon2.rs

//go:noinline
func Poseidon2(state *[16]uint64) {
	ptr := uintptr(unsafe.Pointer(state))
	syscall_poseidon2(ptr)
}

func syscall_poseidon2(ptr uintptr)
