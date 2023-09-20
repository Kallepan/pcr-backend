#!/bin/bash

# Script to use wire to generate a wire_gen.go file aka. dependency injection

# Check if wire is installed
if ! [ -x "$(command -v wire)" ]; then
  echo 'Error: wire is not installed.' >&2
  exit 1
fi

# use wire to generate wire_gen.go
cd src 
wire gen ./config