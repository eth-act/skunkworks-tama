#!/bin/bash
set -e

# Configuration
GOROOT=/home/marcin/projects/lita/EF/skunkworks-tama/tamago-go-latest
WORKDIR=/tmp/manual-build
SRCDIR=.
OUTPUT=empty.elf

# Build tags
TAGS="tamago,linkcpuinit,linkramstart,linkramsize,linkprintk,tinygo.wasm,tinygo,riscv64"

# Memory layout
TEXT_ADDR=0x80001000
DATA_ADDR=0xa0020000

# Environment
export GOROOT=$GOROOT
export GOOS=tamago
export GOARCH=riscv64
export CGO_ENABLED=0

echo "==> Setting up work directory"
rm -rf $WORKDIR
mkdir -p $WORKDIR

echo "==> Creating import configuration for compilation"
cat > $WORKDIR/importcfg.compile << 'EOF'
# Standard library packages (from cache)
packagefile fmt=/home/marcin/.cache/go-build/bc/bc56d912d16b7b4a7244358f30718b84bdf7745af6c6f21f479ebab78caad477-d
packagefile runtime=/home/marcin/.cache/go-build/53/53f7f3fbfcce333583458b16df92a2c727deabc913c9e5880585dae58444bd5b-d
packagefile errors=/home/marcin/.cache/go-build/07/072e7a32aed42f8d4278f0ad898b6317aa91ef591ab50a96afed24750107789d-d
packagefile io=/home/marcin/.cache/go-build/7c/7c673032b9b8eb89ea735335aa579e224ee7035a2dcb348e7d23e4116b07ff1f-d
packagefile math=/home/marcin/.cache/go-build/41/41b4e77f2e6d6a752d99b034d4f2c64ace5885a7a3d156a78959a5600bc92504-d
packagefile os=/home/marcin/.cache/go-build/57/57f9d2f666c56dbdc28e9ca3d86be6eba17771a1bb64bde1e304a1ac97bcecd3-d
packagefile reflect=/home/marcin/.cache/go-build/b0/b0b19cd57f21dd91722c28e4525087c21e2b5164797df464accfa2935283a976-d
packagefile slices=/home/marcin/.cache/go-build/44/444efab5cb0bebe6d1d241e1367d1af691c5b443d3c6fb238a397d0e23097557-d
packagefile strconv=/home/marcin/.cache/go-build/16/167f66596d32c3232d0b7b68638664bd07a76e78510e2a7b78b8ffb5de28fb43-d
packagefile sync=/home/marcin/.cache/go-build/49/498cb503511312baa1ea6f79f72789a151150ae222abb537ca616ca7497ac9a4-d
packagefile tamagotest/tamaboards/zkvm=/home/marcin/.cache/go-build/89/89e6db9b3045fff8ff57dc1b2267fc5d4c212ade363645d35a2a3944685f71eb-d
packagefile github.com/usbarmory/tamago/riscv64=/home/marcin/.cache/go-build/7a/7a195f40aa1e115b31528051a539ef909a6226b20f5217fdd7c9d15dc557a412-d
EOF

echo "==> Compiling source files"
$GOROOT/pkg/tool/linux_amd64/compile \
  -o $WORKDIR/main.o \
  -p main \
  -buildid=manual_build \
  -complete \
  -goversion go1.25.2 \
  -c=4 \
  -nolocalimports \
  -importcfg $WORKDIR/importcfg.compile \
  -pack \
  -d=softfloat \
  $SRCDIR/*.go

$GOROOT/pkg/tool/linux_amd64/asm \
  -I $GOROOT/pkg/include \
  -p main \
  -o $WORKDIR/func_glue.o \
  $SRCDIR/*.s

$GOROOT/pkg/tool/linux_amd64/pack r $WORKDIR/main.a $WORKDIR/main.o $WORKDIR/func_glue.o

echo "==> Creating import configuration for linking"
cat > $WORKDIR/importcfg.link << 'EOF'
packagefile main=$WORKDIR/main.o
packagefile fmt=/home/marcin/.cache/go-build/bc/bc56d912d16b7b4a7244358f30718b84bdf7745af6c6f21f479ebab78caad477-d
packagefile tamagotest/tamaboards/zkvm=/home/marcin/.cache/go-build/89/89e6db9b3045fff8ff57dc1b2267fc5d4c212ade363645d35a2a3944685f71eb-d
packagefile runtime=/home/marcin/.cache/go-build/53/53f7f3fbfcce333583458b16df92a2c727deabc913c9e5880585dae58444bd5b-d
packagefile errors=/home/marcin/.cache/go-build/07/072e7a32aed42f8d4278f0ad898b6317aa91ef591ab50a96afed24750107789d-d
packagefile internal/fmtsort=/home/marcin/.cache/go-build/32/32410364e4bb8b1fd7d3881882612304e711b7e0efb010c9a096273f06ceffbb-d
packagefile io=/home/marcin/.cache/go-build/7c/7c673032b9b8eb89ea735335aa579e224ee7035a2dcb348e7d23e4116b07ff1f-d
packagefile math=/home/marcin/.cache/go-build/41/41b4e77f2e6d6a752d99b034d4f2c64ace5885a7a3d156a78959a5600bc92504-d
packagefile os=/home/marcin/.cache/go-build/57/57f9d2f666c56dbdc28e9ca3d86be6eba17771a1bb64bde1e304a1ac97bcecd3-d
packagefile reflect=/home/marcin/.cache/go-build/b0/b0b19cd57f21dd91722c28e4525087c21e2b5164797df464accfa2935283a976-d
packagefile slices=/home/marcin/.cache/go-build/44/444efab5cb0bebe6d1d241e1367d1af691c5b443d3c6fb238a397d0e23097557-d
packagefile strconv=/home/marcin/.cache/go-build/16/167f66596d32c3232d0b7b68638664bd07a76e78510e2a7b78b8ffb5de28fb43-d
packagefile sync=/home/marcin/.cache/go-build/49/498cb503511312baa1ea6f79f72789a151150ae222abb537ca616ca7497ac9a4-d
packagefile unicode/utf8=/home/marcin/.cache/go-build/82/82fcc1926139add4da78d0fcafb1f56456b5f4d6c7cda324cbb60328dc9e8639-d
packagefile github.com/usbarmory/tamago/riscv64=/home/marcin/.cache/go-build/7a/7a195f40aa1e115b31528051a539ef909a6226b20f5217fdd7c9d15dc557a412-d
packagefile internal/abi=/home/marcin/.cache/go-build/93/936ce508fa0930a6edecdf4839d53b07e4e712c0a90d7ec6a8c2048887283d2b-d
packagefile internal/bytealg=/home/marcin/.cache/go-build/e6/e6b8127fc5c4d9d71a0af9d860ff81f90572436175f4b6e31ce08a07e9c0b17a-d
packagefile internal/byteorder=/home/marcin/.cache/go-build/33/330d8516e1e19124c14ff02305c3e751f4202f29e1e314a0e23a7efba8af533f-d
packagefile internal/chacha8rand=/home/marcin/.cache/go-build/31/31835397461c2cf8784d8e2f87f261ecdd6f15d0c87ff2db1e50dc4eb9cbdab1-d
packagefile internal/coverage/rtcov=/home/marcin/.cache/go-build/92/923c738dbfb7bc1da2cafb1948107f712fc42f7f80fbd6c34645d1e98b35dd82-d
packagefile internal/cpu=/home/marcin/.cache/go-build/0f/0f00c2cd1e80a4eacfaa55ae46ba1ab1bac2e63cb86fb54716e98fa897325b7f-d
packagefile internal/goarch=/home/marcin/.cache/go-build/d6/d65f340f3433765fe0ea3181deafc96f9b25becf72f0def4677dde99a7b493a0-d
packagefile internal/godebugs=/home/marcin/.cache/go-build/b8/b8148a9f64a217c4bde12542d0b4d5ab37e721fbc8b357b5e8519259a8cbf9e3-d
packagefile internal/goexperiment=/home/marcin/.cache/go-build/09/0995ba9c9e4b240dc6e4431eddc0a7d2f62cda3503bc6dfb89fefb6d924f27cb-d
packagefile internal/goos=/home/marcin/.cache/go-build/61/61fb8bf50f22564cb54c9cf8bc70d6e8c26d9bdb7d435551d33ab22c3da5ca3f-d
packagefile internal/profilerecord=/home/marcin/.cache/go-build/96/960e362565073f425ce04884950f9c2ddd50980fd7b5749740828f117782626f-d
packagefile internal/runtime/atomic=/home/marcin/.cache/go-build/b6/b6607c07f149a7999bc46cde6d0a8c48fc250c08f2c8ddae3e9084c44573b3b1-d
packagefile internal/runtime/exithook=/home/marcin/.cache/go-build/b0/b0fdfb9ffeb9189312b880c277006d4c0f26ab65e79fb40477b5ee1df4c2e629-d
packagefile internal/runtime/gc=/home/marcin/.cache/go-build/ab/ab4178291485f28eb110442b381ef8f73b78c0ac63c03dc9b9980bd2caa33edc-d
packagefile internal/runtime/maps=/home/marcin/.cache/go-build/d8/d81c9508eb2d28eed263e7f13ce69332f51d9f17724808bc81e841886ba160dc-d
packagefile internal/runtime/math=/home/marcin/.cache/go-build/f3/f3ef3c8fd9d4fc2f411ca20fdbcacabaaa29f194861a9d715b0f975f84da0708-d
packagefile internal/runtime/strconv=/home/marcin/.cache/go-build/88/88d45ca0bfa75e8ca00a7eab704b10755576f97fbf445af19c4920a94a877599-d
packagefile internal/runtime/sys=/home/marcin/.cache/go-build/60/602b9effab03190bd6b9da1984787a76e506bde411d80f8c9be13bb54c797844-d
packagefile internal/stringslite=/home/marcin/.cache/go-build/ad/ad107e4e1ed3687f2c528ec9219be6dddfb286179193e6c031aaa58098a489f3-d
packagefile internal/trace/tracev2=/home/marcin/.cache/go-build/bf/bfbe89df6b9eae80d78be5585aa265bf7ef7f652cb51f8a77ec3da3ad7d5036c-d
packagefile internal/reflectlite=/home/marcin/.cache/go-build/e6/e6b0857ac35136d385d49f3abceb200c9ac0d73a6076b51c0ed1b428b568daaf-d
packagefile cmp=/home/marcin/.cache/go-build/88/889fda3acc7d8939e025a6619062fc2c6c52c473c7f7e27f5b34ae839ff08e4e-d
packagefile math/bits=/home/marcin/.cache/go-build/82/825f9598368c669ccd3987594913252ef87d0e5c59006858bf8e89239caa304e-d
packagefile internal/filepathlite=/home/marcin/.cache/go-build/a3/a37029c6cc166d936b123856cd3f8a1095617a9217119961ce42d6da9f9d1713-d
packagefile internal/itoa=/home/marcin/.cache/go-build/05/05919c5df3dc92338a758a44cb556938f7efcd1faefe2884be2e7b0f86bf5a03-d
packagefile internal/poll=/home/marcin/.cache/go-build/52/52c53b23225a2969397999ecede254aa48fd7f09c75bef9098aba9135de1a100-d
packagefile internal/syscall/execenv=/home/marcin/.cache/go-build/80/8053e808babf8520577dd136f0dcb034790bc3740d8492a03edd47a9b5cb20dd-d
packagefile internal/syscall/unix=/home/marcin/.cache/go-build/0e/0e43c61df926327c4e198b3c2b270bdfdfd245f62601f645c509cf2438d263b7-d
packagefile internal/testlog=/home/marcin/.cache/go-build/39/39015182b6212c9b9c94dfe4e731556a03690455be2efc330b55783b3fb021cb-d
packagefile io/fs=/home/marcin/.cache/go-build/4a/4a488712952ddd52a703b851e7d88087906bdfba06894848262e091602aea5eb-d
packagefile sync/atomic=/home/marcin/.cache/go-build/53/53e67e5fc826c4449efd11d62acec0417006c56bb1355eca016e2ddcf450f765-d
packagefile syscall=/home/marcin/.cache/go-build/02/02e0be460225997ae3f3b818bdddd12be494d4099b23ed43ef97dfb4583162c1-d
packagefile time=/home/marcin/.cache/go-build/84/843c31be183f698d97bb1715f1177ec5b3bab406a5924dd629a622b0fe41921a-d
packagefile internal/race=/home/marcin/.cache/go-build/53/53bdc443c25abee6f4317c1ebdaaf9dac8e819d4bdf441efc5b1e011bc7ddabd-d
packagefile internal/unsafeheader=/home/marcin/.cache/go-build/86/866c3a271012d21e9dc3c543364a571ccc6c9dea8690f93327156c267afd2e54-d
packagefile iter=/home/marcin/.cache/go-build/b0/b033a0e04e1ae46c8f0277ac9bdb1df3b528262eebf2bf11d216d6227ddb94d9-d
packagefile unicode=/home/marcin/.cache/go-build/89/897dbcd3d5b16fd84023613c77feb915ecb7480b7836b25d17e807760d5d2221-d
packagefile internal/sync=/home/marcin/.cache/go-build/82/82635cca4b9a4da242c6f0b36a6ed53269a9bfb2397b8a4eec8452d35f419bd7-d
packagefile internal/synctest=/home/marcin/.cache/go-build/38/38e9aa3ad5770104ab7c7ac1fb34d2cba2ee289cb3a2255cea182a8a0248520d-d
packagefile github.com/usbarmory/tamago/bits=/home/marcin/.cache/go-build/7e/7e1ef1966335da6425204f7072852bca74c7258bfd0ae39c3a710b527336108b-d
packagefile internal/asan=/home/marcin/.cache/go-build/06/06a84a6342ed2b2248f935d827323375ecd874e0baaa693c0667b04d7a5d165b-d
packagefile internal/msan=/home/marcin/.cache/go-build/1a/1abbbf20105f68ae5fc3b48e54fe2300aa88e28b2bafc2c6bc85d774a3f96987-d
packagefile internal/oserror=/home/marcin/.cache/go-build/a1/a1e2e7807ae60121b673d4928011431e9c898a900c1f338b7766d3460406e80e-d
packagefile path=/home/marcin/.cache/go-build/d3/d3145b1933a3f792df3f056c2ef7993445c3862eff1fdf66015654d90aa291eb-d
packagefile internal/godebug=/home/marcin/.cache/go-build/62/6202416196e25958ffd412324855b61e89dfb0ab72cb70ebfabecdd6ddb82955-d
packagefile internal/bisect=/home/marcin/.cache/go-build/1c/1c0a1dc748d418ae4a726acf5eed4ff7379a2109543955b0877d285b2230d81d-d
EOF

# Replace $WORKDIR placeholder
sed -i "s|\$WORKDIR|$WORKDIR|g" $WORKDIR/importcfg.link

echo "==> Linking with Clang/LLD"
# Method 1: Using go tool link with external linker
$GOROOT/pkg/tool/linux_amd64/link \
  -o $OUTPUT \
  -importcfg $WORKDIR/importcfg.link \
  -buildmode=exe \
  -buildid=manual_build_id \
  -T $TEXT_ADDR \
  -D $DATA_ADDR \
  -extld=clang \
  -extldflags="-target riscv64 -fuse-ld=lld -nostdlib" \
  $WORKDIR/main.a

echo "==> Build complete: $OUTPUT"
ls -lh $OUTPUT

# Optional: Show sections
echo ""
echo "==> ELF sections:"
llvm-readelf -S $OUTPUT | head -20
