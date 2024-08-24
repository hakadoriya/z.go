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
		UsageFunc func(cmd *Command)
		// Description is the description of the command.
		Description string
		// Options is the options of the command.
		Options []Option
		// PreHookFunc is the function to be executed before RunFunc.
		PreHookFunc func(ctx context.Context) error
		// RunFunc is the function to be executed when (*Command).Run is executed.
		RunFunc func(ctx context.Context) error
		// PostHookFunc is the function to be executed after RunFunc.
		PostHookFunc func(ctx context.Context) error
		// SubCommands is the subcommands of the command.
		SubCommands []*Command

		// Stdout is the standard output.
		// If use, get io.Writer from (*Command).GetStdout().
		Stdout io.Writer
		// Stderr is the standard output.
		// If use, get io.Writer from (*Command).GetStderr().
		Stderr io.Writer

		called bool
	}
)

func (cmd *Command) GetCalledCommands() (calledCommands []string) {
	if cmd == nil {
		return nil
	}

	if cmd.called {
		calledCommands = append(calledCommands, cmd.Name)
	}

	for _, subcmd := range cmd.SubCommands {
		commands := subcmd.GetCalledCommands()
		if len(commands) > 0 {
			calledCommands = append(calledCommands, commands...)
		}
	}

	return calledCommands
}

func (cmd *Command) Is(cmdName string) bool {
	if cmd == nil || cmdName == "" {
		return false
	}
	if cmd.Name == cmdName {
		return true
	}
	for _, alias := range cmd.Aliases {
		if alias == cmdName {
			return true
		}
	}
	return false
}

// getSubcommand returns the subcommand if cmd contains the subcommand.
func (cmd *Command) getSubcommand(arg string) (subcmd *Command) {
	if cmd == nil {
		return nil
	}

	for _, subcmd := range cmd.SubCommands {
		if subcmd.Is(arg) {
			return subcmd
		}
	}
	return nil
}
