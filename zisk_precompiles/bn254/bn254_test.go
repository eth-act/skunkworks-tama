package bn254

import (
	_ "tamagotest/tamaboards/zkvm"
	"testing"
)

/*
# Define the bn254 curve parameters
p = 21888242871839275222246405745257275088696311157297823662689037894645226208583
a = 0
b = 3

# Create the elliptic curve
E = EllipticCurve(GF(p), [a, b])

# Define the generator point G
G_x = 1
G_y = 2
G = E(G_x, G_y)

# Add the generator to itself: G + G = 2G
result = G + G + G

# Get the affine coordinates
x_coord = result[0]
y_coord = result[1]

# Print the result
print(f"x = {hex(x_coord)}")
print(f"y = {hex(y_coord)}")
*/

func TestBn254(t *testing.T) {
	p1 := Generator
	p2 := Generator

	params := SyscallBn254AddParams{
		P1: &p1,
		P2: &p2,
	}

	Bn254Dbl(&p1)
	Bn254Add(&params)

	f1 := Complex{
		X: [4]uint64{2, 0, 0, 0},
		Y: [4]uint64{3, 0, 0, 0},
	}
	f2 := Complex{
		X: [4]uint64{4, 0, 0, 0},
		Y: [4]uint64{5, 0, 0, 0},
	}

	Bn254ComplexAdd(&SyscallBn254ComplexAddParams{
		F1: &f1,
		F2: &f2,
	})

	Bn254ComplexMul(&SyscallBn254ComplexMulParams{
		F1: &f1,
		F2: &f2,
	})

	Bn254ComplexSub(&SyscallBn254ComplexSubParams{
		F1: &f1,
		F2: &f2,
	})

	expected_f1 := Complex{
		X: [4]uint64{0x3c208c16d87cfd33, 0x97816a916871ca8d, 0xb85045b68181585d, 0x30644e72e131a029},
		Y: [4]uint64{0x39, 0x0, 0x0, 0x0},
	}

	expected_p1 := Point256{
		X: [4]uint64{0xf2d355961915abf0, 0x9315d84715b8e679, 0xf40232bcb1b6bd15, 0x0769bf9ac56bea3f},
		Y: [4]uint64{0xcdf1ff3dd9fe2261, 0x319e63b40b9c5b57, 0x554fdb7c8d086475, 0x2ab799bee0489429},
	}

	if expected_f1 != f1 || expected_p1 != p1 || p2 != Generator {
		t.Errorf("Failure")
	}
}
