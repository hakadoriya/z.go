package cliz

import (
	"context"
	"strconv"
	"strings"

	"github.com/hakadoriya/z.go/errorz"
)

func (cmd *Command) Parse(ctx context.Context, args []string) (remainingArgs []string, err error) {
	ctx = stdoutWithContext(ctx, Stdout)
	ctx = stderrWithContext(ctx, Stderr)

	remaining, err := cmd.parse(ctx, args)
	if err != nil {
		return nil, errorz.Errorf("%s: cmd.parse: %w", cmd.Name, err)
	}

	return remaining, nil
}

func (cmd *Command) parse(ctx context.Context, args []string) (remainingArgs []string, err error) {
	cmd.initAppendHelpOption()

	if err := cmd.preCheckSubCommands(); err != nil {
		return nil, errorz.Errorf("failed to pre-check commands: %w", err)
	}

	if err := cmd.preCheckOptions(); err != nil {
		return nil, errorz.Errorf("failed to pre-check options: %w", err)
	}

	if err := cmd.loadDefaults(); err != nil {
		return nil, errorz.Errorf("failed to load default: %w", err)
	}

	if err := cmd.loadEnvironments(); err != nil {
		return nil, errorz.Errorf("failed to load environment: %w", err)
	}

	remaining, err := cmd.parseArgs(ctx, args)
	if err != nil {
		return nil, errorz.Errorf("failed to parse arguments: %w", err)
	}

	if err := cmd.checkHelp(ctx); err != nil {
		return nil, err //nolint:wrapcheck
	}

	if err := cmd.postCheckOptions(); err != nil {
		return nil, errorz.Errorf("failed to post-check options: %w", err)
	}

	return remaining, nil
}

func (cmd *Command) parseArgs(ctx context.Context, args []string) (remainingArgs []string, err error) {
	const base, bitSize = 10, 64

	cmd.called = true
	remainingArgs = make([]string, 0)

argsLoop:
	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == breakArg:
			remainingArgs = append(remainingArgs, args[i+1:]...)
			break argsLoop
		case strings.HasPrefix(arg, shortOptionPrefix):
			for _, opt := range cmd.Options {
				switch o := opt.(type) {
				case *StringOption:
					switch {
					case argIsHyphenOption(o, arg):
						if hasNoOptionValue(args, i) {
							return nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
						}
						o.value = ptr(args[i+1])
						i++
						continue argsLoop
					case argIsHyphenOptionEqual(o, arg):
						o.value = ptr(extractValueFromHyphenOptionEqual(arg))
						continue argsLoop
					}
				case *BoolOption:
					switch {
					case argIsHyphenOption(o, arg):
						o.value = ptr(true)
						continue argsLoop
					case argIsHyphenOptionEqual(o, arg):
						optVal, err := strconv.ParseBool(extractValueFromHyphenOptionEqual(arg))
						if err != nil {
							return nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						continue argsLoop
					}
				case *Int64Option:
					switch {
					case argIsHyphenOption(o, arg):
						if hasNoOptionValue(args, i) {
							return nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
						}
						optVal, err := strconv.ParseInt(args[i+1], base, bitSize)
						if err != nil {
							return nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						i++
						continue argsLoop
					case argIsHyphenOptionEqual(o, arg):
						optVal, err := strconv.ParseInt(extractValueFromHyphenOptionEqual(arg), base, bitSize)
						if err != nil {
							return nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						continue argsLoop
					}
				case *Float64Option:
					switch {
					case argIsHyphenOption(o, arg):
						if hasNoOptionValue(args, i) {
							return nil, errorz.Errorf("%s: %w", arg, ErrMissingOptionValue)
						}
						optVal, err := strconv.ParseFloat(args[i+1], bitSize)
						if err != nil {
							return nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						i++
						continue argsLoop
					case argIsHyphenOptionEqual(o, arg):
						optVal, err := strconv.ParseFloat(extractValueFromHyphenOptionEqual(arg), bitSize)
						if err != nil {
							return nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						continue argsLoop
					}
				case *HelpOption:
					switch {
					case argIsHyphenOption(o, arg):
						o.value = ptr(true)
						continue argsLoop
					case argIsHyphenOptionEqual(o, arg):
						optVal, err := strconv.ParseBool(extractValueFromHyphenOptionEqual(arg))
						if err != nil {
							return nil, errorz.Errorf("%s: %w", arg, err)
						}
						o.value = &optVal
						continue argsLoop
					}
				default:
					return nil, errorz.Errorf("%s: %w", arg, ErrInvalidOptionType)
				}
			}
			return nil, errorz.Errorf("%s: %w", arg, ErrUnknownOption)
		default:
			if subcmd := cmd.getSubcommand(arg); subcmd != nil {
				remainingArgs, err = subcmd.parseArgs(ctx, args[i+1:])
				if err != nil {
					return nil, errorz.Errorf("%s: %s: %w", cmd.Name, arg, err)
				}
				return remainingArgs, nil
			}

			// If subcmd is nil, it is not a subcommand.
			remainingArgs = append(remainingArgs, arg)
			continue argsLoop
		}
	}

	return remainingArgs, nil
}
