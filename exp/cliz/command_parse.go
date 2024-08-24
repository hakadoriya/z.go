package cliz

import (
	"context"
	"strconv"
	"strings"

	"github.com/hakadoriya/z.go/errorz"
)

func (cmd *Command) Parse(ctx context.Context, osArgs []string) (remainingArgs []string, err error) {
	remaining, err := cmd.parse(ctx, osArgs)
	if err != nil {
		return nil, errorz.Errorf("%s: cmd.parse: %w", cmd.Name, err)
	}

	return remaining, nil
}

func (cmd *Command) parse(ctx context.Context, osArgs []string) (remainingArgs []string, err error) {
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

	remaining, err := cmd.parseArgs(ctx, osArgs)
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

func (cmd *Command) parseArgs(ctx context.Context, osArgs []string) (remainingArgs []string, err error) {
	const base, bitSize = 10, 64

	cmd.calledCommands = append(cmd.calledCommands, cmd.Name)
	remainingArgs = make([]string, 0)

argsLoop:
	for i := 0; i < len(osArgs); i++ {
		osArg := osArgs[i]

		switch {
		case osArg == breakArg:
			remainingArgs = append(remainingArgs, osArgs[i+1:]...)
			break argsLoop
		case strings.HasPrefix(osArg, shortOptionPrefix):
			for _, opt := range cmd.Options {
				switch o := opt.(type) {
				case *StringOption:
					switch {
					case argIsHyphenOption(o, osArg):
						if hasNoOptionValue(osArgs, i) {
							return nil, errorz.Errorf("%s: %w", osArg, ErrMissingOptionValue)
						}
						o.value = ptr(osArgs[i+1])
						i++
						continue argsLoop
					case argIsHyphenOptionEqual(o, osArg):
						o.value = ptr(extractValueFromHyphenOptionEqual(osArg))
						continue argsLoop
					}
				case *BoolOption:
					switch {
					case argIsHyphenOption(o, osArg):
						o.value = ptr(true)
						continue argsLoop
					case argIsHyphenOptionEqual(o, osArg):
						optVal, err := strconv.ParseBool(extractValueFromHyphenOptionEqual(osArg))
						if err != nil {
							return nil, errorz.Errorf("%s: %w", osArg, err)
						}
						o.value = &optVal
						continue argsLoop
					}
				case *Int64Option:
					switch {
					case argIsHyphenOption(o, osArg):
						if hasNoOptionValue(osArgs, i) {
							return nil, errorz.Errorf("%s: %w", osArg, ErrMissingOptionValue)
						}
						optVal, err := strconv.ParseInt(osArgs[i+1], base, bitSize)
						if err != nil {
							return nil, errorz.Errorf("%s: %w", osArg, err)
						}
						o.value = &optVal
						i++
						continue argsLoop
					case argIsHyphenOptionEqual(o, osArg):
						optVal, err := strconv.ParseInt(extractValueFromHyphenOptionEqual(osArg), base, bitSize)
						if err != nil {
							return nil, errorz.Errorf("%s: %w", osArg, err)
						}
						o.value = &optVal
						continue argsLoop
					}
				case *Float64Option:
					switch {
					case argIsHyphenOption(o, osArg):
						if hasNoOptionValue(osArgs, i) {
							return nil, errorz.Errorf("%s: %w", osArg, ErrMissingOptionValue)
						}
						optVal, err := strconv.ParseFloat(osArgs[i+1], bitSize)
						if err != nil {
							return nil, errorz.Errorf("%s: %w", osArg, err)
						}
						o.value = &optVal
						i++
						continue argsLoop
					case argIsHyphenOptionEqual(o, osArg):
						optVal, err := strconv.ParseFloat(extractValueFromHyphenOptionEqual(osArg), bitSize)
						if err != nil {
							return nil, errorz.Errorf("%s: %w", osArg, err)
						}
						o.value = &optVal
						continue argsLoop
					}
				case *HelpOption:
					switch {
					case argIsHyphenOption(o, osArg):
						o.value = ptr(true)
						continue argsLoop
					case argIsHyphenOptionEqual(o, osArg):
						optVal, err := strconv.ParseBool(extractValueFromHyphenOptionEqual(osArg))
						if err != nil {
							return nil, errorz.Errorf("%s: %w", osArg, err)
						}
						o.value = &optVal
						continue argsLoop
					}
				default:
					return nil, errorz.Errorf("%s: %w", osArg, ErrInvalidOptionType)
				}
			}
			return nil, errorz.Errorf("%s: %w", osArg, ErrUnknownOption)
		default:
			if subcmd := cmd.getSubcommand(osArg); subcmd != nil {
				subcmd.calledCommands = cmd.calledCommands
				defer func() { cmd.calledCommands = subcmd.calledCommands }()
				remainingArgs, err = subcmd.parseArgs(ctx, osArgs[i+1:])
				if err != nil {
					return nil, errorz.Errorf("%s: %s: %w", cmd.Name, osArg, err)
				}
				return remainingArgs, nil
			}

			// If subcmd is nil, it is not a subcommand.
			remainingArgs = append(remainingArgs, osArg)
			continue argsLoop
		}
	}

	return remainingArgs, nil
}
