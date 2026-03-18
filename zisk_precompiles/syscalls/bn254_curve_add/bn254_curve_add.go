//go:build tamago && riscv64

package bn254_curve_add

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/bn254_curve_add.rs

type Point256 struct {
	X [4]uint64
	Y [4]uint64
}

type SyscallBn254AddParams struct {
	P1 *Point256
	P2 *Point256
}

var Generator = Point256{
	X: [4]uint64{0x0000000000000001, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000},
	Y: [4]uint64{0x0000000000000002, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000},
}

var Neutral = Point256{
	X: [4]uint64{0, 0, 0, 0},
	Y: [4]uint64{0, 0, 0, 0},
}

//go:noinline
func Bn254Add(arr *SyscallBn254AddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bn254_add(ptr)
}

func syscall_bn254_add(ptr uintptr)
