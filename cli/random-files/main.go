package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ipfs/go-test/random/files"
)

func main() {
	var usage = `NAME
  %s - Write a random filesystem hierarchy to each <path>

USAGE
  %s [options] <path>...

OPTIONS:
`
	flag.Usage = func() {
		cmd := os.Args[0]
		fmt.Fprintf(os.Stderr, usage, cmd, cmd)
		flag.PrintDefaults()
	}

	var (
		quiet bool
		paths []string
	)

	cfg := files.DefaultConfig()

	flag.BoolVar(&quiet, "q", false, "do not print files and directories")
	flag.IntVar(&cfg.Depth, "depth", cfg.Depth, "depth of the directory tree including the root directory")
	flag.Int64Var(&cfg.FileSize, "filesize", cfg.FileSize, "bytes of random data in each file")
	flag.IntVar(&cfg.Dirs, "dirs", cfg.Dirs, "number of subdirectories at each depth")
	flag.IntVar(&cfg.Files, "files", cfg.Files, "number of files at each depth")
	flag.BoolVar(&cfg.RandomDirs, "random-dirs", cfg.RandomDirs, "randomize number of subdirectories, from 1 to -dirs")
	flag.BoolVar(&cfg.RandomFiles, "random-files", cfg.RandomFiles, "randomize number of files, from 1 to -files")
	flag.BoolVar(&cfg.RandomSize, "random-size", cfg.RandomSize, "randomize file size, from 1 to -filesize")
	flag.Int64Var(&cfg.Seed, "seed", cfg.Seed, "random seed, 0 for current time")
	flag.Parse()

	paths = flag.Args()

	if !quiet {
		cfg.Out = os.Stdout
	}

	err := files.Create(cfg, paths...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		if len(paths) < 1 {
			fmt.Fprintln(os.Stderr)
			flag.Usage()
		}
		os.Exit(1)
	}
}
