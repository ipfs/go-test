package cmd

import "context"

// Watcher is a helper for watching a command's output for a specific string.
// The watcher can watch stdout, stderr or both. It is used by Runner to watch
// for specific output from the commands. The Signal channel is signaled when
// the match string is found.
type Watcher struct {
	match  string
	signal chan struct{}
	stderr bool
	stdout bool
}

// NewWatcher creates a Watcher that is signeled when matching a string from
// stderr or stdout.
func NewWatcher(match string) Watcher {
	return Watcher{
		match:  match,
		signal: make(chan struct{}, 2),
		stderr: true,
		stdout: true,
	}
}

// NewStderrWatcher creates a Watcher that is signeled when matching a string from
// stderr.
func NewStderrWatcher(match string) Watcher {
	return Watcher{
		match:  match,
		signal: make(chan struct{}, 1),
		stderr: true,
	}
}

// NewStdoutWatcher creates a Watcher that is signeled when matching a string from
// stdout.
func NewStdoutWatcher(match string) Watcher {
	return Watcher{
		match:  match,
		signal: make(chan struct{}, 1),
		stdout: true,
	}
}

// Wait waits for the watcher to be signaled for the the context to be
// canceled. If the context is canceled then the context error is returned.
func (w Watcher) Wait(ctx context.Context) error {
	select {
	case <-w.signal:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// Signal returns the channel that is signaled when a line of output matches
// this watcher.
func (w Watcher) Signal() <-chan struct{} {
	return w.signal
}
