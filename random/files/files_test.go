package files_test

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/ipfs/go-test/random/files"
	"github.com/stretchr/testify/require"
)

func TestRandomFiles(t *testing.T) {
	var b bytes.Buffer
	cfg := files.DefaultConfig()
	cfg.Depth = 2
	cfg.Dirs = 5
	cfg.Files = 3
	cfg.Out = &b

	roots := []string{"foo"}
	err := files.Create(cfg, roots...)
	require.NoError(t, err)
	t.Cleanup(func() {
		for _, root := range roots {
			os.RemoveAll(root)
		}
	})

	t.Logf("Created file hierarchy:\n%s", b.String())

	var lines int
	scanner := bufio.NewScanner(&b)
	for scanner.Scan() {
		lines++
	}
	require.NoError(t, scanner.Err())

	subdirs := 0
	if cfg.Depth > 1 {
		dirsAtDepth := cfg.Dirs
		subdirs += dirsAtDepth
		for i := 0; i < cfg.Depth-2; i++ {
			dirsAtDepth *= cfg.Dirs
			subdirs += dirsAtDepth
		}
	}
	linesPerSubdir := cfg.Files + 1
	expect := ((subdirs * linesPerSubdir) + cfg.Files) * len(roots)
	require.Equal(t, expect, lines)
}

func TestRandomFilesValidation(t *testing.T) {
	cfg := files.DefaultConfig()
	err := files.Create(cfg)
	require.Error(t, err)

	cfg.Depth = 0
	require.Error(t, files.Create(cfg, "foo"))
	cfg.Depth = 65
	require.Error(t, files.Create(cfg, "foo"))

	cfg = files.DefaultConfig()

	cfg.Dirs = -1
	require.Error(t, files.Create(cfg, "foo"))
	cfg.Dirs = 65
	require.Error(t, files.Create(cfg, "foo"))

	cfg = files.DefaultConfig()

	cfg.Files = -1
	require.Error(t, files.Create(cfg, "foo"))
	cfg.Files = 65
	require.Error(t, files.Create(cfg, "foo"))

	cfg = files.DefaultConfig()

	cfg.FileSize = -1
	require.Error(t, files.Create(cfg, "foo"))

	cfg = files.DefaultConfig()

	cfg.Depth = 2
	cfg.Dirs = 0
	require.Error(t, files.Create(cfg, "foo"))
}

func TestRandomName(t *testing.T) {
	minSize := 4
	maxSize := 16
	name := files.RandomName(minSize, maxSize)
	require.GreaterOrEqual(t, len(name), minSize)
	require.LessOrEqual(t, len(name), maxSize)
	name = files.RandomName(maxSize, minSize)
	require.GreaterOrEqual(t, len(name), minSize)
	require.LessOrEqual(t, len(name), maxSize)

	fixedSize := 3
	name = files.RandomName(fixedSize)
	require.Len(t, name, fixedSize)

	require.Panics(t, func() {
		files.RandomName(0, maxSize)
	})

	require.Panics(t, func() {
		files.RandomName(0)
	})
}
