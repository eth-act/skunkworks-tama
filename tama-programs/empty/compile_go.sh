#!/bin/bash

CGO_ENABLED=0 GOROOT=/home/marcin/projects/lita/EF/skunkworks-tama/tamago-go-latest GOOS=tamago GOARCH=riscv64 \
  ../../tamago-go-latest/bin/go build \
  -x -v \
  -gcflags="all=-d=softfloat" -ldflags="-T 0x80001000 -D 0xa0020000 -extldflags '-Wl,--whole-archive -L/home/marcin/projects/lita/EF/skunkworks-tama/tama-programs/empty -lfunc -Wl,--no-whole-archive'" \
  -tags tamago,linkcpuinit,linkramstart,linkramsize,linkprintk,tinygo.wasm,tinygo,riscv64 \
  -o empty.elf .

