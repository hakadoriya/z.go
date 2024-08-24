package cliz

import (
	"io"
	"os"
)

func (cmd *Command) GetStdout() io.Writer {
	if cmd.Stdout != nil {
		return cmd.Stdout
	}

	return os.Stdout
}

func (cmd *Command) GetStderr() io.Writer {
	if cmd.Stderr != nil {
		return cmd.Stderr
	}

	return os.Stderr
}
