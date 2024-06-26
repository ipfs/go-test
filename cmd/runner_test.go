package cmd_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ipfs/go-test/cmd"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "out.go")
	err := writeSrc(outPath)
	require.NoError(t, err)

	r := cmd.NewRunner(t, dir)
	werr := cmd.NewStderrWatcher("output on stderr")
	wout := cmd.NewStdoutWatcher("output on stdout")
	wboth := cmd.NewWatcher("output on stdout")
	c := r.Start(context.Background(), cmd.Args("go", "run", outPath), werr, wout, wboth)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = werr.Wait(ctx)
	require.NoError(t, err)
	t.Log("stderr watcher signaled")

	err = wout.Wait(ctx)
	require.NoError(t, err)
	t.Log("stdout watcher signaled")

	err = wboth.Wait(ctx)
	require.NoError(t, err)
	t.Log("both watcher signaled")

	r.Stop(c, time.Second)
}

var outSrc = `
package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	fmt.Println("output on stdout")
	fmt.Fprintln(os.Stderr, "output on stderr")

	<-shutdown
}
`

func writeSrc(outPath string) error {
	outGo, err := os.Create(outPath)
	if err != nil {
		return err
	}
	_, err = outGo.WriteString(outSrc)
	if err != nil {
		outGo.Close()
		return err
	}
	return outGo.Close()
}
