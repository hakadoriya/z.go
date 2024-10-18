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
		PreHookExecFunc func(c *Command, args []string) error
		// ExecFunc is the function to be executed when (*Command).Exec is executed.
		ExecFunc func(c *Command, args []string) error
		// PostHookExecFunc is the function to be executed after ExecFunc.
		PostHookExecFunc func(c *Command, args []string) error
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

		// ctx is the context.
		// MEMO: If the structure does not include context.Context, and you want to use the value
		//       set by WithValue in PreHookExecFunc in ExecFunc, PreHookExecFunc must return context.Context,
		//       but ExecFunc and PostHookExecFunc are also in the same relationship.
		//       However, it is inconvenient to have to return ctx with return considering the subsequent processing,
		//       so I decided to include context.Context in this structure.
		// MEMO: もし context.Context を含まない場合, PreHookExecFunc で WithValue した値を
		//       ExecFunc で利用したい場合, PreHookExecFunc は context.Context を返す必要があるが,
		//       同様に ExecFunc と PostHookExecFunc も同じ関係にある.
		//       しかし, 後続の処理の事を考えて return で ctx を返さなければならないのは不便であるため,
		//       この構造体に context.Context を含めることにした.
		//
		//nolint:containedctx
		ctx                     context.Context
		allExecutedCommandNames []string
		remainingArgs           []string
	}
)

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
