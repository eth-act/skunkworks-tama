//go:build tamago && riscv64

package sha256

import (
	"unsafe"
)

type SyscallSha256Params struct {
	State *[4]uint64
	Input *[8]uint64
}

//go:noinline
func Sha256(arr *SyscallSha256Params) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_sha256(ptr)
}

func syscall_sha256(ptr uintptr)
