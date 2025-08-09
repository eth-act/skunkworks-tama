//go:build tamago && riscv64

package main

import (
	"fmt"
	_ "tamagotest/tamaboards/zkvm"
)

func main() {
	a := 42
	b := 7
	
	fmt.Printf("Calculator Test Program\n")
	fmt.Printf("a = %d, b = %d\n", a, b)
	fmt.Printf("Addition: %d + %d = %d\n", a, b, a+b)
	fmt.Printf("Subtraction: %d - %d = %d\n", a, b, a-b)
	fmt.Printf("Multiplication: %d * %d = %d\n", a, b, a*b)
	fmt.Printf("Division: %d / %d = %d\n", a, b, a/b)
	fmt.Printf("Modulo: %d %% %d = %d\n", a, b, a%b)
	
	c := 15
	d := 4
	fmt.Printf("\nBitwise operations with c = %d, d = %d\n", c, d)
	fmt.Printf("AND: %d & %d = %d\n", c, d, c&d)
	fmt.Printf("OR: %d | %d = %d\n", c, d, c|d)
	fmt.Printf("XOR: %d ^ %d = %d\n", c, d, c^d)
	fmt.Printf("Left shift: %d << %d = %d\n", c, 2, c<<2)
	fmt.Printf("Right shift: %d >> %d = %d\n", c, 1, c>>1)
}