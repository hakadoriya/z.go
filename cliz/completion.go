package cliz

import (
	"text/template"

	"github.com/hakadoriya/z.go/cliz/clicorez"
	"github.com/hakadoriya/z.go/errorz"
)

type completionBashTmplData struct {
	RootCommandName                      string
	GenerateBashCompletionSubCommandName string
}

func (c *Command) initAppendCompletionSubCommand() {
	c.SubCommands = append(
		c.SubCommands,
		//nolint:exhaustruct
		&Command{
			Name:   clicorez.DefaultCompletionSubCommandName,
			Hidden: true,
			SubCommands: []*Command{
				{
					Name:        "bash",
					Description: "generate bash completion script",
					ExecFunc: func(cmd *Command, _ []string) error {
						b, err := completionBashTmplFile.ReadFile(completionBashTmpl)
						if err != nil {
							return errorz.Errorf("completionBashTmplFile.ReadFile: name=%s: %w", completionBashTmpl, err)
						}
						tmpl := template.Must(template.New(completionBashTmpl).Parse(string(b)))
						return tmpl.Execute(cmd.Stdout(), completionBashTmplData{
							RootCommandName:                      c.Name,
							GenerateBashCompletionSubCommandName: clicorez.DefaultGenerateBashCompletionSubCommandName,
						})
					},
				},
			},
		})
}
