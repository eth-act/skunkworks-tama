//go:build tamago && riscv64

package sha256f

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/sha256f.rs

type SyscallSha256Params struct {
	State *[4]uint64
	Input *[8]uint64
}

//go:noinline
func Sha256Update(arr *SyscallSha256Params) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_sha256f(ptr)
}

func syscall_sha256f(ptr uintptr)
