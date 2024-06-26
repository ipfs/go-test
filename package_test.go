package test

import (
	"testing"

	"github.com/ipfs/go-test/cmd"
	"github.com/ipfs/go-test/random"
)

func TestImports(t *testing.T) {
	args := cmd.Args("echo", "hello", "world")
	if args.String() != "echo hello world" {
		t.Error("something is wrong with cmd")
	}

	const naddrs = 3
	addrs := random.Addrs(naddrs)
	if len(addrs) != naddrs {
		t.Error("something is wrong random")
	}
}
