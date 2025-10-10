FROM ubuntu:22.04

# Required by ZISK
RUN apt update && apt install -y libclang-dev xz-utils jq curl build-essential qemu-system libomp-dev libgmp-dev nlohmann-json3-dev protobuf-compiler uuid-dev libgrpc++-dev libsecp256k1-dev libsodium-dev libpqxx-dev nasm libopenmpi-dev openmpi-bin openmpi-common

# Utils
RUN apt update && apt install -y git vim curl wget

RUN git config --global --add safe.directory '*'

# Install rust
RUN curl https://sh.rustup.rs -sSf | bash -s -- -y

# Install go
RUN wget https://go.dev/dl/go1.25.2.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.25.2.linux-amd64.tar.gz && rm go1.25.2.linux-amd64.tar.gz

ENV PATH="$PATH:/usr/local/go/bin"

