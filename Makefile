.PHONY: all clean build-tamago build-zisk compile-empty

TAMAGO_DIR = tamago-go-latest
TAMAGO_SRC = $(TAMAGO_DIR)/src
TAMAGO_BIN = $(TAMAGO_DIR)/bin
TAMAGO = $(TAMAGO_BIN)/go

ZISK_DIR = zisk
ZISKEMU = $(ZISK_DIR)/target/release/ziskemu

# Compilation flags for TamaGo
GCFLAGS = -gcflags="all=-d=softfloat"
LDFLAGS = -ldflags="-T 0x80000000"
TAGS = -tags tamago,linkcpuinit,linkramstart,linkramsize,linkprintk

all: build-tamago build-zisk

build-tamago:
	cd $(TAMAGO_SRC) && ./all.bash
	@cd $(TAMAGO_BIN) && export TAMAGO=$$(pwd)/go && echo "TAMAGO set to: $$TAMAGO"

build-zisk:
	cd $(ZISK_DIR) && cargo build --release -p ziskemu
	@cd $(ZISK_DIR) && export ZISKEMU=$$(pwd)/target/release/ziskemu && echo "ZISKEMU set to: $$ZISKEMU"

build: build-tamago build-zisk

clean:
	rm -rf $(TAMAGO_DIR)
	rm -f latest.zip
	cd $(ZISK_DIR) && cargo clean

compile-empty:
	cd tama-programs/empty && GOOS=tamago GOARCH=riscv64 ../../$(TAMAGO) build $(GCFLAGS) $(LDFLAGS) $(TAGS) -o empty.elf .

run-empty: compile-empty
	cd tama-programs/empty && ../../$(ZISKEMU) -e empty.elf -i empty_input.bin

trace-empty: compile-empty
	cd tama-programs/empty && ../../$(ZISKEMU) -e empty.elf -i empty_input.bin -a -v

log-empty: compile-empty
	cd tama-programs/empty && ../../$(ZISKEMU) -e empty.elf -i empty_input.bin -l -p 1

trace-file-empty: compile-empty
	cd tama-programs/empty && ../../$(ZISKEMU) -e empty.elf -i empty_input.bin -a -t trace.out && echo "Trace saved to tama-programs/empty/trace.out"