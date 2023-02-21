package fio

import (
	"bytes"
	"testing"
)

func TestNewCommandWrapper(t *testing.T) {
	cmd := newCommandWrapper("echo", "hello world")
	if cmd == nil {
		t.Error("Expected CommandWrapper, got nil")
	}
}

func TestCommandWrapper_Run(t *testing.T) {
	cmd := newCommandWrapper("echo", "hello world")

	output, err := cmd.Run()
	if err != nil {
		t.Errorf("Expected nil error, got %s", err)
	}
	if !bytes.Equal(output, []byte("hello world\n")) {
		t.Errorf("Expected 'hello world\\n', got %s", string(output))
	}
}

func TestCommandWrapper_Run_Error(t *testing.T) {
	cmd := newCommandWrapper("invalid_command", "hello world")

	_, err := cmd.Run()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
