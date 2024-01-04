#!/bin/bash

cd src
go mod download

# install go tools
go install github.com/google/wire/cmd/wire@latest