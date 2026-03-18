package twist

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Expected values captured from ziskemu v0.16.0

func TestDblLine(t *testing.T) {
	p := [24]uint64{
		// x_c0
		0xf5f28fa202940a10, 0xb3f5fb2687b4961a, 0xa1a893b53e2ae580, 0x9894999d1a3caee9,
		0x6f67b7631863366b, 0x058191924350bcd7,
		// x_c1
		0xa5a9c0759e23f606, 0xaaa0c59dbccd60c3, 0x3bb17e18e2867806, 0x1b1ab6cc8541b367,
		0xc2b6ed0ef2158547, 0x11922a097360edf3,
		// y_c0
		0x4c730af860494c4a, 0x597cfa1f5e369c5a, 0xe7e6856caa0a635a, 0xbbefb5e96e0d495f,
		0x07d3a975f0ef25a2, 0x0083fd8e7e80dae5,
		// y_c1
		0xadc0fc92df64b05d, 0x18aa270a2b1461dc, 0x86adac6a3be4eba0, 0x79495c4ec93da33a,
		0xe7175850a43ccaed, 0x0b2bc2a163de1bf2,
	}

	coeffs := FcallDblLineCoeffs(&p)

	expectedLambda := [12]uint64{
		0x17b7da8790b423fd, 0x61654e65baad75fa, 0xfd4f2a7852ba4aa0, 0x28b1d389469afb43,
		0x4e9f7858be382288, 0x1104cf27d1bcba46, 0x11a8dc6587e093db, 0x67199f1b6216eb02,
		0xcf8b4ba361c58d3e, 0xdfca408f771dd978, 0x143f283c1fc43484, 0x0f60d32ddd934c88,
	}
	expectedMu := [12]uint64{
		0xa5e1288731199f5a, 0x59c932cf3e947040, 0x14d0950c0b15fdaf, 0xec2f40d8ccd55fd9,
		0x051bb2c04d587e69, 0x0c8699b080019d87, 0x70031c402cfaf3f6, 0xcf3b6ee318e5f54a,
		0xdb6019aa9ba12cc7, 0x17d4d2894abfd387, 0xbd1deb4a967336f8, 0x06547f814c53fbde,
	}

	if coeffs.Lambda != expectedLambda {
		t.Errorf("lambda: got %v, want %v", coeffs.Lambda, expectedLambda)
	}
	if coeffs.Mu != expectedMu {
		t.Errorf("mu: got %v, want %v", coeffs.Mu, expectedMu)
	}
}
