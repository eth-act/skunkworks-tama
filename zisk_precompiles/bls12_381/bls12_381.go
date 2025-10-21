//go:build tamago && riscv64

package bls12_381

import (
	"unsafe"
)

type Point256 struct {
	X [6]uint64
	Y [6]uint64
}

type Complex struct {
	X [6]uint64
	Y [6]uint64
}

type SyscallBls12_381ComplexAddParams struct {
	F1 *Complex
	F2 *Complex
}

type SyscallBls12_381ComplexMulParams struct {
	F1 *Complex
	F2 *Complex
}

type SyscallBls12_381ComplexSubParams struct {
	F1 *Complex
	F2 *Complex
}

type SyscallBls12_381AddParams struct {
	P1 *Point256
	P2 *Point256
}

var Generator = Point256{
	X: [6]uint64{0xfb3af00adb22c6bb, 0x6c55e83ff97a1aef, 0xa14e3a3f171bac58, 0xc3688c4f9774b905, 0x2695638c4fa9ac0f, 0x17f1d3a73197d794},
	Y: [6]uint64{0x0caa232946c5e7e1, 0xd03cc744a2888ae4, 0x00db18cb2c04b3ed, 0xfcf5e095d5d00af6, 0xa09e30ed741d8ae4, 0x08b3f481e3aaa0f1},
}

var Neutral = Point256{
	X: [6]uint64{0, 0, 0, 0, 0, 0},
	Y: [6]uint64{0, 0, 0, 0, 0, 0},
}

//go:noinline
func Bls12_381Add(arr *SyscallBls12_381AddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bls12_381_add(ptr)
}

//go:noinline
func Bls12_381Dbl(arr *Point256) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bls12_381_dbl(ptr)
}

//go:noinline
func Bls12_381ComplexAdd(arr *SyscallBls12_381ComplexAddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bls12_381_complex_add(ptr)
}

//go:noinline
func Bls12_381ComplexMul(arr *SyscallBls12_381ComplexMulParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bls12_381_complex_mul(ptr)
}

//go:noinline
func Bls12_381ComplexSub(arr *SyscallBls12_381ComplexSubParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bls12_381_complex_sub(ptr)
}

func syscall_bls12_381_add(ptr uintptr)
func syscall_bls12_381_dbl(ptr uintptr)
func syscall_bls12_381_complex_add(ptr uintptr)
func syscall_bls12_381_complex_mul(ptr uintptr)
func syscall_bls12_381_complex_sub(ptr uintptr)
