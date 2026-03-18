//go:build tamago && riscv64

package secp256k1_add

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/secp256k1_add.rs

type Point256 struct {
	X [4]uint64
	Y [4]uint64
}

type SyscallSecp256k1AddParams struct {
	P1 *Point256
	P2 *Point256
}

var Generator = Point256{
	X: [4]uint64{0x59F2815B16F81798, 0x029BFCDB2DCE28D9, 0x55A06295CE870B07, 0x79BE667EF9DCBBAC},
	Y: [4]uint64{0x9C47D08FFB10D4B8, 0xFD17B448A6855419, 0x5DA4FBFC0E1108A8, 0x483ADA7726A3C465},
}

var Neutral = Point256{
	X: [4]uint64{0, 0, 0, 0},
	Y: [4]uint64{0, 0, 0, 0},
}

// Given points `p1` and `p2`, performs the point addition `p1 + p2` and assigns the result to `p1`.
// It assumes that `p1` and `p2` are from the Secp256k1 curve, that `p1,p2 != 𝒪` and that `p2 != p1,-p1`

//go:noinline
func Secp256k1Add(arr *SyscallSecp256k1AddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_secp256k1_add(ptr)
}

func syscall_secp256k1_add(ptr uintptr)
