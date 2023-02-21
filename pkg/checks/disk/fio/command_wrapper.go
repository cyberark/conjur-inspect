package fio

import (
	"bytes"
	"os/exec"
)

type command interface {
	Run() ([]byte, error)
}

type commandWrapper struct {
	command *exec.Cmd
	stdout  bytes.Buffer
	stderr  bytes.Buffer
}

func newCommandWrapper(name string, args ...string) command {
	wrapper := commandWrapper{}

	// Instantiate the command
	command := exec.Command(name, args...)

	// Bind the stdout and stderr
	command.Stdout = &wrapper.stdout
	command.Stderr = &wrapper.stderr

	// Wrap the command
	wrapper.command = command

	return &wrapper
}

func (wrapper *commandWrapper) Run() ([]byte, error) {
	err := wrapper.command.Run()

	if err != nil {
		return nil, err
	}

	return wrapper.stdout.Bytes(), nil
}
