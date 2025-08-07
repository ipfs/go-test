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
	MinimumNameSize = 4
	fileNameAlpha   = "abcdefghijklmnopqrstuvwxyz01234567890-_"
)

// Config contains settings for creating random files and directories.
type Config struct {
	// Depth is the depth of the directory tree including the root directory.
	Depth int
	// Dirs is the number of subdirectories at each depth.
	Dirs int
	// Files is the number of files at each depth.
	Files int
	// FileSize sets the number of random bytes in each file.
	FileSize int64
	// NameMaxSize is the maximum length of a random file or directory name.
	NameMaxSize int
	// NameMinSize is the minimum length of a random file or directory name. It
	// must be at least MinimumNameSize.
	NameMinSize int
	// Where to write display output, such as os.Stdout. Default is nil.
	Out io.Writer
	// RandomDirss specifies whether or not to randomize the number of
	// subdirectories from 1 to the value configured by Dirs.
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

// DefaultConfig returns default settings for creating random files and
// directories.
func DefaultConfig() Config {
	return Config{
		Depth:       2,
		Dirs:        5,
		Files:       10,
		FileSize:    4096,
		NameMaxSize: 16,
		NameMinSize: 4,
		RandomDirs:  false,
		RandomFiles: false,
		RandomSize:  true,
	}
}

// Create creates random files and directories according to the provided
// configuration. The random files and directories are created in the specified
// root paths.
func Create(cfg Config, roots ...string) error {
	if len(roots) == 0 {
		return errors.New("must provide at least 1 root directory path")
	}
	err := cfg.validate()
	if err != nil {
		return err
	}

	var rnd *rand.Rand
	if cfg.Seed == 0 {
		rnd = random.NewRand()
	} else {
		rnd = random.NewSeededRand(cfg.Seed)
	}

	for _, root := range roots {
		err := os.MkdirAll(root, 0755)
		if err != nil {
			return err
		}

		err = cfg.writeTree(rnd, root, 1)
		if err != nil {
			return err
		}

	}

	return nil
}

// RandomName generates a random file or directory name.
//
// If no sizes are specified, then the default minimum and maximum name sizes
// are used. If one size is specified, then the name will have that size. If
// two sizes are specified, then the name will have a random size between the
// smaller and larger of the two numbers.
func RandomName(sizes ...int) string {
	var cfg Config
	if len(sizes) > 0 {
		var minSize, maxSize int
		if len(sizes) > 1 {
			// Random size between minimum and maximum.
			minSize = min(sizes[0], sizes[1])
			maxSize = max(sizes[0], sizes[1])
		} else {
			// Fixes size as specified.
			minSize = sizes[0]
			maxSize = minSize
		}
		err := validateNameSize(minSize, maxSize)
		if err != nil {
			panic(err)
		}
		cfg.NameMaxSize = maxSize
		cfg.NameMinSize = minSize
	} else {
		// Use default random size.
		cfg = DefaultConfig()
	}

	return cfg.randomName(random.NewRand())
}

func (cfg *Config) validate() error {
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
	err := validateNameSize(cfg.NameMinSize, cfg.NameMaxSize)
	if err != nil {
		return err
	}

	return nil
}

func validateNameSize(minSize, maxSize int) error {
	if minSize < MinimumNameSize {
		return fmt.Errorf("minimum name size must be at least %d", MinimumNameSize)
	}
	if maxSize < minSize {
		return errors.New("maximum name size is less than minimum name size")
	}
	return nil
}

func (cfg *Config) writeTree(rnd *rand.Rand, root string, depth int) error {
	nFiles := cfg.Files
	if nFiles != 0 {
		if cfg.RandomFiles && nFiles > 1 {
			nFiles = rnd.Intn(nFiles) + 1
		}

		for i := 0; i < nFiles; i++ {
			if err := cfg.writeFile(rnd, root); err != nil {
				return err
			}
		}
	}

	return cfg.writeSubdirs(rnd, root, depth)
}

func (cfg *Config) writeSubdirs(rnd *rand.Rand, root string, depth int) error {
	if depth == cfg.Depth {
		return nil
	}
	depth++

	nDirs := cfg.Dirs
	if cfg.RandomDirs && nDirs > 1 {
		nDirs = rnd.Intn(nDirs) + 1
	}

	for i := 0; i < nDirs; i++ {
		if err := cfg.writeSubdir(rnd, root, depth); err != nil {
			return err
		}
	}

	return nil
}

func (cfg *Config) writeSubdir(rnd *rand.Rand, root string, depth int) error {
	name := cfg.randomName(rnd)
	root = filepath.Join(root, name)
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}

	if cfg.Out != nil {
		fmt.Fprintln(cfg.Out, root+"/")
	}

	return cfg.writeTree(rnd, root, depth)
}

func (cfg *Config) randomName(rnd *rand.Rand) string {
	sizeDiff := cfg.NameMaxSize - cfg.NameMinSize
	n := cfg.NameMinSize
	if sizeDiff != 0 {
		n += rnd.Intn(sizeDiff)
	}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = fileNameAlpha[rnd.Intn(len(fileNameAlpha))]
	}
	return string(b)
}

func (cfg *Config) writeFile(rnd *rand.Rand, root string) error {
	name := cfg.randomName(rnd)
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
