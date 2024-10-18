package cliz

import (
	"io"
)

func (c *Command) Stdout() io.Writer {
	if c.stdout != nil {
		return c.stdout
	}

	return Stdout
}

func (c *Command) Stderr() io.Writer {
	if c.stderr != nil {
		return c.stderr
	}

	return Stderr
}

func (c *Command) SetStdout(w io.Writer) {
	if c == nil {
		return
	}

	c.stdout = w
}

func (c *Command) SetStderr(w io.Writer) {
	if c == nil {
		return
	}

	c.stderr = w
}

func (c *Command) SetStdoutRecursive(w io.Writer) {
	if c == nil {
		return
	}

	c.SetStdout(w)
	for _, subcmd := range c.SubCommands {
		subcmd.SetStdoutRecursive(w)
	}
}

func (c *Command) SetStderrRecursive(w io.Writer) {
	if c == nil {
		return
	}

	c.SetStderr(w)
	for _, subcmd := range c.SubCommands {
		subcmd.SetStderrRecursive(w)
	}
}
