//go:build tamago && riscv64

package main

import (
	"fmt"
	_ "tamagotest/tamaboards/zkvm"
)

// Follow instruction from: https://words.filippo.io/rustgo/

// myExternalFunction is implemented in external_riscv64.s
// It bridges to a C function that takes two uint32 and returns uint32
// go:noescape
func myExternalFunction(a, b uint32) uint32

func main() {	
	x := uint32(10)
        y := uint32(20)
	result := myExternalFunction(x, y)
	fmt.Printf("myExternalFunction(%d, %d) = %d\n", x, y, result)
}
