package main

import (
	"io"
	"os"

	"github.com/conjurinc/conjur-preflight/pkg/cmd"
)

var cmdStdout, cmdStderr io.Writer

func init() {
	cmdStdout = os.Stdout
	cmdStderr = os.Stderr
}

func main() {
	cmd.Execute(cmdStdout, cmdStderr)
}
