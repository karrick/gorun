# gorun

Small Go library to invoke child proceses and wait for their completion

Documentation is available via
[![GoDoc](https://godoc.org/github.com/karrick/gorun?status.svg)](https://godoc.org/github.com/karrick/gorun)
and
[https://pkg.go.dev/github.com/karrick/gorun?tab=doc](https://pkg.go.dev/github.com/karrick/gorun?tab=doc).

## Description

Executes a system command.

1. When this cannot spawn the requested program, it returns a nil
Response and an error.

2. When it can spawn the requested program, but receives a system
call error while waiting for the child program to exit, it returns
an Response with a -1 exit Code, and an Err set to system call
error message.

3. When the child program exited due to receiving a signal, it
returns an Response with a -1 exit Code -1 and an error message
that specifies the signal it receives in Err.

NOTE: If context.Context expires, Go will send termination signal
to spawned child process, and Response will have Code and Err set
in accordance with this case.

4. When the child program exits on its own and not due to receiving
a signal as described above, it returns Response with Code set
to the exit code of the child program, and Err set to nil.

## Example

```Go
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/karrick/gorun"
)

func main() {
	req := &gorun.Request{
		Path:  "/usr/bin/cat",
		Stdin: []byte("line 1\nline 2\n"),
	}
	resp, err := req.Run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
	fmt.Printf("%q Err:\t%v\n", req.Path, resp.Err)
	fmt.Printf("%q Code:\t%v\n", req.Path, resp.Code)
	fmt.Printf("%q Stdout:\t%q\n", req.Path, string(resp.Stdout))
	fmt.Printf("%q Stderr:\t%q\n", req.Path, string(resp.Stderr))
}
```
