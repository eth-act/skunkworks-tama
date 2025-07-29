.PHONY: all clean build-tamago build-zisk

TAMAGO_DIR = tamago-go-latest
TAMAGO_SRC = $(TAMAGO_DIR)/src
TAMAGO_BIN = $(TAMAGO_DIR)/bin

ZISK_DIR = zisk

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

run-empty:
	cd tama-programs/empty && $(ZISK_DIR)/target/release/ziskemu -e empty.elf -i empty_input.bin