//go:build !windows
// +build !windows

package gorun

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Run("cannot spawn", func(t *testing.T) {
		t.Run("path empty string", func(t *testing.T) {
			_, err := Run(context.Background(), &Request{})
			ensureError(t, err, ErrSpawn{Err: errors.New("fork/exec : no such file or directory")})
		})
		t.Run("no such executable", func(t *testing.T) {
			_, err := Run(context.Background(), &Request{
				Path: "/no-such-path",
			})
			ensureError(t, err, ErrSpawn{Err: errors.New("fork/exec /no-such-path: no such file or directory")})
		})
		t.Run("invalid dir", func(t *testing.T) {
			_, err := Run(context.Background(), &Request{
				Path: "/usr/bin/true",
				Dir:  "/does-not-exist",
			})
			ensureError(t, err, ErrSpawn{Err: errors.New("chdir /does-not-exist: no such file or directory")})
		})
	})

	t.Run("true", func(t *testing.T) {
		got, err := Run(context.Background(), &Request{
			Path: "/usr/bin/true",
		})
		ensureError(t, err, nil)
		want := &Response{
			Stderr: []byte{},
			Stdout: []byte{},
		}
		ensureResponsesMatch(t, got, want)
	})

	t.Run("false", func(t *testing.T) {
		got, err := Run(context.Background(), &Request{
			Path: "/usr/bin/false",
		})
		ensureError(t, err, nil)
		want := &Response{
			Code:   1,
			Stderr: []byte{},
			Stdout: []byte{},
		}
		ensureResponsesMatch(t, got, want)
	})
	t.Run("echo", func(t *testing.T) {
		got, err := Run(context.Background(), &Request{
			Path: "/bin/echo",
			Args: []string{"one", "two", "three"},
		})
		ensureError(t, err, nil)
		want := &Response{
			Stderr: []byte{},
			Stdout: []byte("one two three\n"),
		}
		ensureResponsesMatch(t, got, want)
	})
	t.Run("test-script.sh", func(t *testing.T) {
		got, err := Run(context.Background(), &Request{
			Path:  "./test-script.sh",
			Env:   []string{"GORUN=asdf"},
			Args:  []string{"first", "line", "from", "arguments"},
			Stdin: strings.NewReader("second line from stdin"),
		})
		ensureError(t, err, nil)
		want := &Response{
			Code:   13,
			Stderr: []byte("prints to standard error: asdf\n"),
			Stdout: []byte("first line from arguments\nsecond line from stdin"),
		}
		ensureResponsesMatch(t, got, want)
	})
	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		got, err := Run(ctx, &Request{
			Path: "/bin/sleep",
			Args: []string{"1"},
		})
		ensureError(t, err, nil)
		want := &Response{
			Code:   -1,
			Err:    ErrSignal{Err: errors.New("signal: killed")},
			Stderr: []byte{},
			Stdout: []byte{},
		}
		ensureResponsesMatch(t, got, want)
		if !errors.Is(got.Err, want.Err) {
			t.Errorf("GOT: %T(%v); WANT: %T(%v)", got.Err, got.Err, want.Err, want.Err)
		}
	})
	t.Run("canceled", func(t *testing.T) {
		t.Run("before start", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := Run(ctx, &Request{
				Path: "/bin/sleep",
				Args: []string{"1"},
			})

			ensureError(t, err, ErrSpawn{Err: errors.New("context canceled")})
		})
		t.Run("after start", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			go func() {
				time.Sleep(100 * time.Millisecond)
				cancel()
			}()

			got, err := Run(ctx, &Request{
				Path: "/bin/sleep",
				Args: []string{"1"},
			})

			ensureError(t, err, nil)

			want := &Response{
				Code:   -1,
				Err:    ErrSignal{Err: errors.New("signal: killed")},
				Stderr: []byte{},
				Stdout: []byte{},
			}

			ensureResponsesMatch(t, got, want)
		})
	})
}
