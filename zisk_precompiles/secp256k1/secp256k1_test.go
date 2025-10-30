package secp256k1

import (
	_ "tamagotest/tamaboards/zkvm"
	"testing"
)

/*
# Define the secp256k1 curve parameters
p = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F
a = 0
b = 7

# Create the elliptic curve
E = EllipticCurve(GF(p), [a, b])

# Define the generator point G
G_x = 0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798
G_y = 0x483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8
G = E(G_x, G_y)

result = G + G + G

# Get the affine coordinates
x_coord = result[0]
y_coord = result[1]

# Print the result
print(f"x = {hex(x_coord)}")
print(f"y = {hex(y_coord)}")
*/

func TestSecp256k1AddDbl(t *testing.T) {
	p1 := Generator
	p2 := Generator

	params := SyscallSecp256k1AddParams{
		P1: &p1,
		P2: &p2,
	}

	Secp256k1Dbl(&p1)
	Secp256k1Add(&params)

	expected_p1 := Point256{
		X: [4]uint64{0x8601f113bce036f9, 0xb531c845836f99b0, 0x49344f85f89d5229, 0xf9308a019258c310},
		Y: [4]uint64{0x6cb9fd7584b8e672, 0x6500a99934c2231b, 0x0fe337e62a37f356, 0x388f7b0f632de814},
	}

	if expected_p1 != p1 || p2 != Generator {
		t.Errorf("Failure")
	}
}
