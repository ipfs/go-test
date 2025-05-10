# :test_tube: go-test
[![Go Reference](https://pkg.go.dev/badge/github.com/ipfs/go-test.svg)](https://pkg.go.dev/github.com/ipfs/go-test)
[![Go Test](https://github.com/ipfs/go-test/actions/workflows/go-test.yml/badge.svg)](https://github.com/ipfs/go-test/actions/workflows/go-test.yml)
> A testing utility library.

This module provides testing utility logic in Go that is not specific to any project.

## [`cmd`](https://pkg.go.dev/github.com/ipfs/go-test/cmd "API documentation") package

The cmd package contains logic for running synchronous and asynchronous commands.

## [`random`](https://pkg.go.dev/github.com/ipfs/go-test/random "API documentation") package

The random package contains logic for generating random test data.

## Command Line Tools

Command line utilities are located in the [`cli`](https://github.com/ipfs/go-test/tree/main/cli) directory:
- [random-data](https://github.com/ipfs/go-test/tree/main/cli/random-data#random-data---writes-random-data-to-stdout) writes random data to stdout
- [random-files](https://github.com/ipfs/go-test/tree/main/cli/random-files#random-files---create-random-filesystem-hierarchies) creates random files in hierarchy of random directories
