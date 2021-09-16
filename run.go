package gorun

import (
	"bytes"
	"context"
	"io"
	"os/exec"
)

// Run executes a system command.
//
// 1. When this cannot spawn the requested program, it returns a nil
// Response and an error.
//
// 2. When it can spawn the requested program, but receives a system
// call error while waiting for the child program to exit, it returns
// an Response with a -1 exit Code, and an Err set to system call
// error message.
//
// 3. When the child program exited due to receiving a signal, it
// returns an Response with a -1 exit Code -1 and an error message
// that specifies the signal it receives in Err.
//
// NOTE: If context.Context expires, Go will send termination signal
// to spawned child process, and Response will have Code and Err set
// in accordance with this case.
//
// 4. When the child program exits on its own and not due to receiving
// a signal as described above, it returns Response with Code set
// to the exit code of the child program, and Err set to nil.
func Run(ctx context.Context, req *Request) (*Response, error) {
	// NOTE: Calling this function gets optimized out by the Go
	// compiler and transformed into the method invocation.
	return req.Run(ctx)
}

// Request represents a request to spawn a child process.
type Request struct {
	// Args is a potentially empty list of command line arguments to
	// be sent to the child process.
	Args []string

	// Env is a potentially empty list of environment variable
	// assignments to be sent to the child process.
	Env []string

	// Stdin is the potentially nil io.Reader that will be available
	// for the child process to read from when it reads from its
	// standard input.
	Stdin io.Reader

	// Dir is the directory to set as the child process' initial
	// current working directory when it starts.
	Dir string

	// Path is the path to the child process program executable file.
	Path string
}

// Run executes a system command.
//
// 1. When this cannot spawn the requested program, it returns a nil
// Response and an error.
//
// 2. When it can spawn the requested program, but receives a system
// call error while waiting for the child program to exit, it returns
// an Response with a -1 exit Code, and an Err set to system call
// error message.
//
// 3. When the child program exited due to receiving a signal, it
// returns an Response with a -1 exit Code -1 and an error message
// that specifies the signal it receives in Err.
//
// NOTE: If context.Context expires, Go will send termination signal
// to spawned child process, and Response will have Code and Err set
// in accordance with this case.
//
// 4. When the child program exits on its own and not due to receiving
// a signal as described above, it returns Response with Code set
// to the exit code of the child program, and Err set to nil.
func (req *Request) Run(ctx context.Context) (*Response, error) {
	var stderr, stdout bytes.Buffer
	var err error

	cmd := exec.CommandContext(ctx, req.Path, req.Args...)
	cmd.Dir = req.Dir
	cmd.Env = req.Env
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if req.Stdin != nil {
		cmd.Stdin = req.Stdin
	}

	if err = cmd.Start(); err != nil {
		return nil, ErrSpawn{Err: err}
	}

	err = cmd.Wait()

	resp := &Response{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
	}

	// Go standard library interprets whether a child program was
	// successful based on its exit code. However many programs this
	// expects to invoke work properly and return information in the
	// exit code. Because of the difference in interpretation of the
	// exit code, namely, a non-zero exit code does not imply program
	// failure, this needs to handle the case when program terminated
	// due to a signal, and call that an error, but a non-zero exit
	// code not due to a signal is not itself an error.
	switch e := err.(type) {
	case nil:
		// happy case: note Code is already 0 which is exit code of program
		return resp, nil
	case exitCoder:
		// Go standard library returns an error that implements
		// exitCoder when the child process either returns a non-zero
		// exit code, which includes when the process terminates from
		// a signal.
		resp.Code = e.ExitCode()
		if resp.Code == -1 {
			// Go standard library returns exit code -1 when program
			// has either not yet exited, or when it was terminated by
			// a signal. Because this library only checks the exit
			// code after the child program exits, it is only -1 when
			// the child program exited due to receiving a signal.
			resp.Err = ErrSignal{Err: err}
		}
		return resp, nil
	default:
		// Some other meta error due to trying to manage child
		// process.
		return nil, ErrWait{Err: err}
	}
}

type exitCoder interface {
	ExitCode() int
}

// Response represents the result of spawning a child process.
type Response struct {
	// Err will be nil when it was able to spawn the child process and
	// the child process terminated without being sent a signal. Err
	// will be non-nil when the program could not properly spawn or
	// collect the exit status of the child process, or when the child
	// process exited as a result of receiving a signal.
	Err error

	// Stderr will be a potentially empty slice of bytes that
	// represent whatever the child process wrote to its standard
	// error file stream.
	Stderr []byte

	// Stdout will be a potentially empty slice of bytes that
	// represent whatever the child process wrote to its standard
	// output file stream.
	Stdout []byte

	// Code will be the exit code that the child process returned when
	// it exited. When its value is -1, the child process was spawned
	// but terminated in response to receiving a signal.
	Code int
}

type ErrSignal struct {
	Err error
}

func (e ErrSignal) Error() string {
	return e.Err.Error()
}

func (e ErrSignal) Is(err error) bool {
	_, ok := err.(ErrSignal)
	return ok
}

func (e ErrSignal) Unwrap() error { return e.Err }

type ErrSpawn struct {
	Err error
}

func (e ErrSpawn) Error() string {
	return "cannot spawn process: " + e.Err.Error()
}

func (e ErrSpawn) Is(err error) bool {
	_, ok := err.(ErrSpawn)
	return ok
}

func (e ErrSpawn) Unwrap() error { return e.Err }

type ErrWait struct {
	Err error
}

func (e ErrWait) Error() string {
	return "cannot wait for process to terminate: " + e.Err.Error()
}

func (e ErrWait) Is(err error) bool {
	_, ok := err.(ErrWait)
	return ok
}

func (e ErrWait) Unwrap() error { return e.Err }
