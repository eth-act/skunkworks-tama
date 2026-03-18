//go:build tamago && riscv64

package input

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/input.rs

//go:noinline
func FcallInput(addr *uint64) {
	fcall_input(uintptr(*addr))
}

func fcall_input(addr uintptr)
