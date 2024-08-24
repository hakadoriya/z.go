package cliz

import (
	"embed"
	"io"
	"strings"
)

//nolint:gochecknoglobals
var (
	//go:embed completion_bash.tmpl
	completionBashTmplFile embed.FS
	completionBashTmpl     = "completion_bash.tmpl"
)

func (c *Command) initAppendGenerateBashCompletionSubCommands() {
	// recursively
	for _, subcmd := range c.SubCommands {
		subcmd.initAppendGenerateBashCompletionSubCommands()
	}

	c.SubCommands = append(
		c.SubCommands,
		//nolint:exhaustruct
		&Command{
			Name:   DefaultGenerateBashCompletionSubCommandName,
			Hidden: true,
			ExecFunc: func(_ *Command, _ []string) error {
				completions := make([]string, 0)
				for _, subcmd := range c.SubCommands {
					if subcmd.Hidden {
						continue
					}

					completions = append(completions, subcmd.Name)
					completions = append(completions, subcmd.Aliases...)
				}

				for _, option := range c.Options {
					completions = append(completions, longOptionPrefix+option.GetName())
					for _, alias := range option.GetAliases() {
						completions = append(completions, shortOptionPrefix+alias)
					}
				}

				_, _ = io.WriteString(c.Stdout(), strings.Join(completions, " ")+"\n")
				return nil
			},
		})
}
