package gorun

import (
	"errors"
	"strings"
	"testing"
)

var someError = errors.New("some error message")

func ensureError(tb testing.TB, got, want error) {
	tb.Helper()
	if got == nil {
		if want != nil {
			tb.Fatalf("GOT: %v; WANT: %T(%q)", got, want, want.Error())
		}
	} else if want == nil {
		tb.Fatalf("GOT: %T(%q); WANT: %v", got, got.Error(), want)
	} else {
		var target error
		if ok := errors.As(got, &target); !ok {
			tb.Fatalf("GOT: %T(%q); WANT: %T(%q)", got, got.Error(), want, want.Error())
		}
		if g, w := got.Error(), want.Error(); !strings.Contains(g, w) {
			tb.Fatalf("GOT: %v; WANT: %v", g, w)
		}
	}
}

func TestEnsureError(t *testing.T) {
	t.Run("non nil", func(t *testing.T) {
		got := someError
		want := someError
		ensureError(t, got, want)
	})
	t.Run("nil", func(t *testing.T) {
		var got, want error
		ensureError(t, got, want)
	})
}

func ensureResponsesMatch(tb testing.TB, got, want *Response) {
	tb.Helper()
	if got == nil {
		if want != nil {
			tb.Fatalf("GOT: %v; WANT: %v", got, want)
		}
	} else if want == nil {
		tb.Fatalf("GOT: %v; WANT: %v", got, want)
	} else {
		if got.Code != want.Code {
			tb.Fatalf("GOT: %v; WANT: %v", got.Code, want.Code)
		}
		ensureError(tb, got.Err, want.Err)
		if g, w := string(got.Stderr), string(want.Stderr); g != w {
			tb.Fatalf("GOT: %q; WANT: %q", g, w)
		}
		if g, w := string(got.Stdout), string(want.Stdout); g != w {
			tb.Fatalf("GOT: %q; WANT: %q", g, w)
		}
	}
}

func TestEnsureResponsesMatch(t *testing.T) {
	got := &Response{
		Err:    someError,
		Code:   13,
		Stderr: []byte("error 1\nerror 2\n"),
		Stdout: []byte("out 1\nout 2\n"),
	}
	want := &Response{
		Err:    someError,
		Code:   13,
		Stderr: []byte("error 1\nerror 2\n"),
		Stdout: []byte("out 1\nout 2\n"),
	}
	ensureResponsesMatch(t, got, want)
}
