package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hakadoriya/z.go/cliz"
	"github.com/hakadoriya/z.go/errorz"
)

func main() {
	c := cliz.Command{
		Name: os.Args[0],
		ExecFunc: func(c *cliz.Command, args []string) error {
			const o = `This is a example command.
Tips: Exec below commands for auto completion:
	eval "$(./example completion bash)"
`
			_, _ = io.WriteString(c.Stdout(), o)
			return nil
		},
		SubCommands: []*cliz.Command{
			{
				Name: "good-morning",
				SubCommands: []*cliz.Command{
					{
						Name:        "world",
						Description: "Prints 'Good morning, world!'",
						Options: []cliz.Option{
							&cliz.StringOption{
								Name:        "who",
								Description: "If you want to say 'Good morning, who!', set this option.",
							},
						},
						ExecFunc: func(c *cliz.Command, _ []string) error {
							who, err := c.GetOptionString("who")
							if err != nil {
								return errorz.Errorf("c.GetOptionString: %w", err)
							}
							if who != "" {
								fmt.Fprintf(c.Stdout(), "Good morning, %s!\n", who)
							} else {
								fmt.Fprintln(c.Stdout(), "Good morning, world!")
							}
							return nil
						},
					},
				},
			},
			{
				Name: "hello",
				SubCommands: []*cliz.Command{
					{
						Name:        "world",
						Description: "Prints 'Hello, world!'",
						Options: []cliz.Option{
							&cliz.StringOption{
								Name:        "who",
								Description: "If you want to say 'Good morning, who!', set this option.",
							},
						},
						ExecFunc: func(c *cliz.Command, _ []string) error {
							who, err := c.GetOptionString("who")
							if err != nil {
								return errorz.Errorf("c.GetOptionString: %w", err)
							}
							if who != "" {
								fmt.Fprintf(c.Stdout(), "Hello, %s!\n", who)
							} else {
								fmt.Fprintln(c.Stdout(), "Hello, world!")
							}
							return nil
						},
					},
				},
			},
			{
				Name: "good-night",
				SubCommands: []*cliz.Command{
					{
						Name:        "world",
						Description: "Prints 'Good night, world!'",
						Options: []cliz.Option{
							&cliz.StringOption{
								Name:        "who",
								Description: "If you want to say 'Good morning, who!', set this option.",
							},
						},
						ExecFunc: func(c *cliz.Command, args []string) error {
							who, err := c.GetOptionString("who")
							if err != nil {
								return errorz.Errorf("c.GetOptionString: %w", err)
							}
							if who != "" {
								fmt.Fprintf(c.Stdout(), "Good night, %s!\n", who)
							} else {
								fmt.Fprintln(c.Stdout(), "Good night, world!")
							}
							return nil
						},
					},
				},
			},
		},
	}

	if err := c.Exec(context.Background(), os.Args); err != nil {
		if cliz.IsHelp(err) {
			os.Exit(1)
		}
		err = errorz.Errorf("c.Run: %w", err)
		log.Fatalf("%+v", err)
	}
}
