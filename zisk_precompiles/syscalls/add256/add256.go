//go:build tamago && riscv64

package add256

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/add256.rs

type SyscallAdd256Params struct {
	A   *[4]uint64
	B   *[4]uint64
	Cin uint64
	C   *[4]uint64
}

//go:noinline
func Add256(params *SyscallAdd256Params) uint64 {
	ptr := uintptr(unsafe.Pointer(params))
	return syscall_add256(ptr)
}

func syscall_add256(ptr uintptr) uint64
