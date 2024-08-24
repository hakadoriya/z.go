package cliz

import (
	"context"
	"io"
)

type (
	// Command is a structure for building command lines. Please fill in each field for the structure you are facing.
	Command struct {
		// Name is the name of the command.
		Name string
		// Aliases is the alias names of the command.
		Aliases []string
		// Group is the group name of the command.
		Group string
		// Usage is the usage of the command.
		//
		// If you want to use the default usage, remain empty.
		// Otherwise, set the custom usage.
		Usage string
		// UsageFunc is custom usage function.
		//
		// If you want to use the default usage function, remain nil.
		// Otherwise, set the custom usage function.
		UsageFunc func(c *Command)
		// Description is the description of the command.
		Description string
		// Options is the options of the command.
		Options []Option
		// PreHookExecFunc is the function to be executed before ExecFunc.
		PreHookExecFunc func(ctx context.Context, rootCmd Cmd, args []string) error
		// ExecFunc is the function to be executed when (*Command).Exec is executed.
		ExecFunc func(ctx context.Context, rootCmd Cmd, args []string) error
		// PostHookExecFunc is the function to be executed after ExecFunc.
		PostHookExecFunc func(ctx context.Context, rootCmd Cmd, args []string) error
		// SubCommands is the subcommands of the command.
		SubCommands []*Command
		// If Hidden is true, the command is not displayed in the help message and completion.
		Hidden bool

		// stdout is the standard output.
		// If use in ExecFunc, get io.Writer from (command).GetStdout().
		stdout io.Writer
		// stderr is the standard output.
		// If use in ExecFunc, get io.Writer from (command).GetStderr().
		stderr io.Writer

		allExecutedCommandNames []string
	}

	Cmd interface {
		GetExecutedCommand() *Command
		GetExecutedCommandNames() []string
		GetStdout() io.Writer
		GetStderr() io.Writer
		GetOptionString(name string) (string, error)
		GetOptionBool(name string) (bool, error)
		GetOptionInt64(name string) (int64, error)
		GetOptionFloat64(name string) (float64, error)
		private() bool
	}
)

func (c *Command) private() bool {
	return true
}

func (c *Command) GetExecutedCommand() *Command {
	if c == nil {
		return nil
	}

	if len(c.allExecutedCommandNames) > 0 {
		for _, subcmd := range c.SubCommands {
			if len(subcmd.allExecutedCommandNames) > 0 {
				return subcmd.GetExecutedCommand()
			}
		}

		return c
	}

	return nil
}

func (c *Command) GetExecutedCommandNames() (calledCommands []string) {
	if c == nil {
		return nil
	}

	return c.allExecutedCommandNames
}

func (c *Command) is(cmdName string) bool {
	if c == nil || cmdName == "" {
		return false
	}
	if c.Name == cmdName {
		return true
	}
	for _, alias := range c.Aliases {
		if alias == cmdName {
			return true
		}
	}
	return false
}

// getSubcommand returns the subcommand if cmd contains the subcommand.
func (c *Command) getSubcommand(arg string) (subcmd *Command) {
	if c == nil {
		return nil
	}

	for _, subcmd := range c.SubCommands {
		if subcmd.is(arg) {
			return subcmd
		}
	}
	return nil
}
