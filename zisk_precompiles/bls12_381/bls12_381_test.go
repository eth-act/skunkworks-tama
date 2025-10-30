package bls12_381

import (
	_ "tamagotest/tamaboards/zkvm"
	"testing"
)

/*
# Define the bls12_381 curve parameters
p = 0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab
a = 0
b = 4

# Create the elliptic curve
E = EllipticCurve(GF(p), [a, b])

# Define the generator point G
G_x = 0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb
G_y = 0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1
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

func TestBls12381(t *testing.T) {
	p1 := Generator
	p2 := Generator

	params := SyscallBls12_381AddParams{
		P1: &p1,
		P2: &p2,
	}

	Bls12_381Dbl(&p1)
	Bls12_381Add(&params)

	f1 := Complex{
		X: [6]uint64{2, 0, 0, 0, 0, 0},
		Y: [6]uint64{3, 0, 0, 0, 0, 0},
	}
	f2 := Complex{
		X: [6]uint64{4, 0, 0, 0, 0, 0},
		Y: [6]uint64{5, 0, 0, 0, 0, 0},
	}

	Bls12_381ComplexAdd(&SyscallBls12_381ComplexAddParams{
		F1: &f1,
		F2: &f2,
	})

	Bls12_381ComplexMul(&SyscallBls12_381ComplexMulParams{
		F1: &f1,
		F2: &f2,
	})

	Bls12_381ComplexSub(&SyscallBls12_381ComplexSubParams{
		F1: &f1,
		F2: &f2,
	})

	expected_f1 := Complex{
		X: [6]uint64{0xb9feffffffffaa97, 0x1eabfffeb153ffff, 0x6730d2a0f6b0f624, 0x64774b84f38512bf, 0x4b1ba7b6434bacd7, 0x1a0111ea397fe69a},
		Y: [6]uint64{0x39, 0x0, 0x0, 0x0, 0x0, 0x0},
	}

	// that was obtained with sagemath
	expected_p1 := Point256{
		X: [6]uint64{0x96d2c0c9024e5224, 0x81747a0b2ca2179b, 0xf3780a51335b3ff9, 0xb112d61f9be9a5f1, 0x1765212deca99697, 0x09ece308f9d1f013},
		Y: [6]uint64{0xa3473b0590ae30d1, 0xe745256c634af45c, 0x9d9c27310fd43be6, 0xa69a0cddabc3097f, 0x8a84623389c5f80c, 0x032b80d3a6f5b09f},
	}

	if expected_f1 != f1 || p1 != expected_p1 || p2 != Generator {
		t.Errorf("Failure")
	}
}
