package cliz

import (
	"context"
	"strconv"
	"strings"

	"github.com/hakadoriya/z.go/contextz"
	"github.com/hakadoriya/z.go/errorz"
)

func (c *Command) parse(ctx context.Context, osArgs []string) (remainingArgs []string, err error) {
	if len(c.remainingArgs) > 0 {
		// If remainingArgs is not empty, it is already parsed.
		return c.remainingArgs, nil
	}

	c.ctx = ctx

	// following is not idempotent.

	c.initAppendHelpOption()
	c.initAppendCompletionSubCommand()
	c.initAppendGenerateBashCompletionSubCommands()

	if err := c.preCheckSubCommands(); err != nil {
		return nil, errorz.Errorf("failed to pre-check commands: %w", err)
	}

	if err := c.preCheckOptions(); err != nil {
		return nil, errorz.Errorf("failed to pre-check options: %w", err)
	}

	if err := c.loadDefaults(); err != nil {
		return nil, errorz.Errorf("failed to load default: %w", err)
	}

	if err := c.loadEnvironments(); err != nil {
		return nil, errorz.Errorf("failed to load environment: %w", err)
	}

	remaining, err := c.parseArgs(osArgs)
	if err != nil {
		return nil, errorz.Errorf("failed to parse arguments: %w", err)
	}

	if err := c.checkHelp(); err != nil {
		return nil, err
	}

	if err := c.postCheckOptions(); err != nil {
		return nil, errorz.Errorf("failed to post-check options: %w", err)
	}

	if err := contextz.CheckContext(ctx); err != nil {
		return nil, errorz.Errorf("failed to check context: %w", err)
	}

	return remaining, nil
}

func ptr[T interface{}](v T) *T { return &v }

//nolint:cyclop,funlen,gocognit,gocyclo
func (c *Command) parseArgs(osArgs []string) (remainingArgs []string, err error) {
	defer func() { c.remainingArgs = remainingArgs }()
	const base, bitSize = 10, 64

	c.allExecutedCommandNames = append(c.allExecutedCommandNames, c.Name)
	remainingArgs = make([]string, 0)

argsLoop:
	for i := 0; i < len(osArgs); i++ {
		osArg := osArgs[i]

		switch {
		case osArg == breakArg:
			remainingArgs = append(remainingArgs, osArgs[i+1:]...)
			break argsLoop
		case strings.HasPrefix(osArg, shortOptionPrefix):
			for _, opt := range c.Options {
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
				case *Uint64Option:
					switch {
					case argIsHyphenOption(o, osArg):
						if hasNoOptionValue(osArgs, i) {
							return nil, errorz.Errorf("%s: %w", osArg, ErrMissingOptionValue)
						}
						optVal, err := strconv.ParseUint(osArgs[i+1], base, bitSize)
						if err != nil {
							return nil, errorz.Errorf("%s: %w", osArg, err)
						}
						o.value = &optVal
						i++
						continue argsLoop
					case argIsHyphenOptionEqual(o, osArg):
						optVal, err := strconv.ParseUint(extractValueFromHyphenOptionEqual(osArg), base, bitSize)
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
			if subcmd := c.getSubcommand(osArg); subcmd != nil {
				//nolint:fatcontext
				subcmd.ctx = c.ctx
				subcmd.allExecutedCommandNames = c.allExecutedCommandNames
				defer func() {
					// NOTE: Propagate the updated subcmd.ctx to c.ctx in subcmd processing.
					//
					//nolint:fatcontext
					c.ctx = subcmd.ctx
					c.allExecutedCommandNames = subcmd.allExecutedCommandNames
				}()
				remainingArgs, err = subcmd.parseArgs(osArgs[i+1:])
				if err != nil {
					return nil, errorz.Errorf("%s: %s: %w", c.Name, osArg, err)
				}
				return remainingArgs, nil
			}

			// If sub is nil, it is not a subcommand.
			remainingArgs = append(remainingArgs, osArg)
			continue argsLoop
		}
	}

	return remainingArgs, nil
}
