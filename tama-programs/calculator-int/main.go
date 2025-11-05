//go:build tamago && riscv64

package main

import (
	"fmt"
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"github.com/eth-act/skunkworks-tama/tamaboards/zkvm/zisk_runtime"
)

func main() {
	// Read two integers using the generic Read function
	a := zisk_runtime.Read[int64]()
	b := zisk_runtime.Read[int64]()

	fmt.Printf("Calculator Test Program\n")

	fmt.Printf("a = %d, b = %d\n", a, b)

	fmt.Printf("Addition: %d + %d = %d\n", a, b, a+b)
	fmt.Printf("Subtraction: %d - %d = %d\n", a, b, a-b)
	fmt.Printf("Multiplication: %d * %d = %d\n", a, b, a*b)
	fmt.Printf("Division: %d / %d = %d\n", a, b, a/b)
	fmt.Printf("Modulo: %d %% %d = %d\n", a, b, a%b)

	fmt.Printf("AND: %d & %d = %d\n", a, b, a&b)
	fmt.Printf("OR: %d | %d = %d\n", a, b, a|b)
	fmt.Printf("XOR: %d ^ %d = %d\n", a, b, a^b)
	fmt.Printf("Left shift: %d << %d = %d\n", a, 2, a<<2)
	fmt.Printf("Right shift: %d >> %d = %d\n", a, 1, a>>1)
}
