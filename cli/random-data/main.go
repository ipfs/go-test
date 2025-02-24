package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"

	random "github.com/ipfs/go-test/random"
)

func main() {
	var usage = `NAME
  %s - Write random data to stdout

USAGE
  %s [options]

OPTIONS:
`
	flag.Usage = func() {
		cmd := os.Args[0]
		fmt.Fprintf(os.Stderr, usage, cmd, cmd)
		flag.PrintDefaults()
	}

	var (
		b64  bool
		seed int64
		size int64
	)
	flag.BoolVar(&b64, "b64", false, "base-64 encode output")
	flag.Int64Var(&seed, "seed", 0, "random seed, 0 or unset for current time")
	flag.Int64Var(&size, "size", 0, "number of bytes to generate")
	flag.Parse()

	if size < 1 {
		fmt.Fprintln(os.Stderr, "missing value for size")
		fmt.Fprintln(os.Stderr)
		flag.Usage()
		os.Exit(1)
	}

	err := writeData(seed, size, b64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "missing value for size")
		os.Exit(1)
	}
	if b64 {
		fmt.Println()
	}
}

func writeData(seed, size int64, b64 bool) error {
	var rnd *rand.Rand
	if seed == 0 {
		rnd = random.NewRand()
	} else {
		rnd = random.NewSeededRand(seed)
	}

	var w io.Writer
	if b64 {
		b := base64.NewEncoder(base64.StdEncoding, os.Stdout)
		defer b.Close()
		w = b
	} else {
		b := bufio.NewWriter(os.Stdout)
		defer b.Flush()
		w = b
	}
	_, err := io.CopyN(w, rnd, size)
	return err
}
