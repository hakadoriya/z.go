package cliz

import (
	"io"
	"os"
)

func (cmd *Command) GetStdout() io.Writer {
	if cmd.stdout != nil {
		return cmd.stdout
	}

	return os.Stdout
}

func (cmd *Command) GetStderr() io.Writer {
	if cmd.stderr != nil {
		return cmd.stderr
	}

	return os.Stderr
}

func (cmd *Command) SetStdout(w io.Writer) {
	if cmd == nil {
		return
	}

	cmd.stdout = w
}

func (cmd *Command) SetStderr(w io.Writer) {
	if cmd == nil {
		return
	}

	cmd.stderr = w
}

func (cmd *Command) SetStdoutRecursive(w io.Writer) {
	if cmd == nil {
		return
	}

	cmd.SetStdout(w)
	for _, subcmd := range cmd.SubCommands {
		subcmd.SetStdoutRecursive(w)
	}
}

func (cmd *Command) SetStderrRecursive(w io.Writer) {
	if cmd == nil {
		return
	}

	cmd.SetStderr(w)
	for _, subcmd := range cmd.SubCommands {
		subcmd.SetStderrRecursive(w)
	}
}
