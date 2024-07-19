#!/bin/bash

# Set the desired architecture
ARCH="amd64"

# Run go generate
go generate ./...

# Build for the specified architecture
GOOS=linux GOARCH=$ARCH go build -o build/linux_$ARCH/moba-converter-go

echo "Build complete for Linux ($ARCH)"