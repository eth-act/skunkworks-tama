#!/bin/bash

set -e

if [ -z "$ZISKEMU" ]; then
    echo "Error: ZISKEMU environment variable is not set"
    exit 1
fi

$ZISKEMU -c --elf $1

