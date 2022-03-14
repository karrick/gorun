package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/karrick/gorun"
)

func main() {
	show(&gorun.Request{
		Path: "/bin/dates",
	})

	show(&gorun.Request{
		Path: "/usr/bin/true",
	})

	show(&gorun.Request{
		Path: "/usr/bin/false",
	})

	show(&gorun.Request{
		Path: "/bin/sleep",
		Args: []string{"300"},
	})

	show(&gorun.Request{
		Path:  "/bin/cat",
		Stdin: strings.NewReader("line 1\nline 2\n"),
	})

	show(&gorun.Request{
		Path:  "/bin/echo",
		Stdin: strings.NewReader("line 1\nline 2\n"),
		Args:  []string{"one", "two"},
	})
}

func show(req *gorun.Request) {
	resp, err := req.Run(context.Background())
	if err != nil {
		fmt.Printf("\n%q invocation err:\t%v\n", req.Path, err)
		return
	}
	fmt.Printf("\n%q Err:\t%v\n", req.Path, resp.Err)
	fmt.Printf("%q Code:\t%v\n", req.Path, resp.Code)
	fmt.Printf("%q Stdout:\t%q\n", req.Path, string(resp.Stdout))
	fmt.Printf("%q Stderr:\t%q\n", req.Path, string(resp.Stderr))
}
