package files

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/ipfs/go-test/random"
)

const (
	fileNameSize  = 16
	fileNameAlpha = "abcdefghijklmnopqrstuvwxyz01234567890-_"
)

type Config struct {
	// Depth is the depth of the directory tree including the root directory.
	Depth int
	// Dirs is the number of subdirectories at each depth.
	Dirs int
	// Files is the number of files at each depth.
	Files int
	// FileSize sets the number of random bytes in each file.
	FileSize int64
	// Where to write display output, such as os.Stdout. Default is nil.
	Out io.Writer
	// RandomDirss specifies whether or not to randomize the number of
	// subdirectoriess from 1 to the value configured by Dirs.
	RandomDirs bool
	// RandomFiles specifies whether or not to randomize the number of files
	// from 1 to the value configured by Files.
	RandomFiles bool
	// RandomSize specifies whether or not to randomize the file size from 1 to
	// the value configured by FileSize.
	RandomSize bool
	// Seed sets the seen for the random number generator when set to a
	// non-zero value.
	Seed int64
}

func DefaultConfig() Config {
	return Config{
		Depth:       2,
		Dirs:        5,
		Files:       10,
		FileSize:    4096,
		RandomDirs:  false,
		RandomFiles: false,
		RandomSize:  true,
	}
}

func validateConfig(cfg *Config) error {
	if cfg.Depth < 1 || cfg.Depth > 64 {
		return errors.New("depth out of range, must be between 1 and 64")
	}
	if cfg.Dirs < 0 || cfg.Dirs > 64 {
		return errors.New("dirs out of range, must be between 0 and 64")
	}
	if cfg.Files < 0 || cfg.Files > 64 {
		return errors.New("files out of range, must be between 0 and 64")
	}
	if cfg.FileSize < 0 {
		return errors.New("file size out of range, must be 0 or greater")
	}
	if cfg.Depth > 1 && cfg.Dirs < 1 {
		return errors.New("dirs must be at least 1 for depth > 1")
	}

	return nil
}

func Create(cfg Config, paths ...string) error {
	if len(paths) == 0 {
		return errors.New("must provide at least 1 root directory path")
	}
	err := validateConfig(&cfg)
	if err != nil {
		return err
	}

	var rnd *rand.Rand
	if cfg.Seed == 0 {
		rnd = random.NewRand()
	} else {
		rnd = random.NewSeededRand(cfg.Seed)
	}

	for _, root := range paths {
		err := os.MkdirAll(root, 0755)
		if err != nil {
			return err
		}

		err = writeTree(rnd, root, 1, &cfg)
		if err != nil {
			return err
		}

	}

	return nil
}

func writeTree(rnd *rand.Rand, root string, depth int, cfg *Config) error {
	nFiles := cfg.Files
	if nFiles != 0 {
		if cfg.RandomFiles && nFiles > 1 {
			nFiles = rnd.Intn(nFiles) + 1
		}

		for i := 0; i < nFiles; i++ {
			if err := writeFile(rnd, root, cfg); err != nil {
				return err
			}
		}
	}

	return writeSubdirs(rnd, root, depth, cfg)
}

func writeSubdirs(rnd *rand.Rand, root string, depth int, cfg *Config) error {
	if depth == cfg.Depth {
		return nil
	}
	depth++

	nDirs := cfg.Dirs
	if cfg.RandomDirs && nDirs > 1 {
		nDirs = rnd.Intn(nDirs) + 1
	}

	for i := 0; i < nDirs; i++ {
		if err := writeSubdir(rnd, root, depth, cfg); err != nil {
			return err
		}
	}

	return nil
}

func writeSubdir(rnd *rand.Rand, root string, depth int, cfg *Config) error {
	name := randomFilename(rnd)
	root = filepath.Join(root, name)
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}

	if cfg.Out != nil {
		fmt.Fprintln(cfg.Out, root+"/")
	}

	return writeTree(rnd, root, depth, cfg)
}

func randomFilename(rnd *rand.Rand) string {
	n := rnd.Intn(fileNameSize-4) + 4
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = fileNameAlpha[rnd.Intn(len(fileNameAlpha))]
	}
	return string(b)
}

func writeFile(rnd *rand.Rand, root string, cfg *Config) error {
	name := randomFilename(rnd)
	filePath := filepath.Join(root, name)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	if cfg.FileSize > 0 {
		fileSize := cfg.FileSize
		if cfg.RandomSize && fileSize > 1 {
			fileSize = rnd.Int63n(fileSize) + 1
		}

		if _, err := io.CopyN(f, rnd, fileSize); err != nil {
			f.Close()
			return err
		}
	}

	if cfg.Out != nil {
		fmt.Fprintln(cfg.Out, filePath)
	}

	return f.Close()
}
