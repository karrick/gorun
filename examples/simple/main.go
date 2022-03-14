package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/karrick/gorun"
)

func main() {
	req := &gorun.Request{
		Path:  "/bin/cat",
		Stdin: strings.NewReader("line 1\nline 2\n"),
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
