#!/bin/bash

# This script is used to launch the development environment for the

export $(grep -v '^#' ./.dev.env | xargs)

cd src
go run main.go
