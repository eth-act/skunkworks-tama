#!/bin/bash

set -e

pushd .
cd docker
docker build -t skunkworks-tama .
popd

docker run \
  --rm -it \
  --volume .:/skunkworks-tama \
  --workdir /skunkworks-tama \
  skunkworks-tama

