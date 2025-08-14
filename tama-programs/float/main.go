//go:build tamago && riscv64

package main

import (
	"fmt"
	_ "tamagotest/tamaboards/zkvm"
)

func main() {

	var f1 float64 = 3.14
	var f2 float64 = 2.71
	var result float64 = f1 + f2

	fmt.Printf("Float operation completed: %f + %f = %f\n", f1, f2, result)
}
