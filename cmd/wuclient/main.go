package main

import (
	"fmt"
	"os"

	"github.com/mas9612/wrapups/internal/pkg/command"
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("wuclient", "v0.2.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"list": func() (cli.Command, error) {
			return &command.ListCommand{}, nil
		},
		"get": func() (cli.Command, error) {
			return &command.GetCommand{}, nil
		},
		"create": func() (cli.Command, error) {
			return &command.CreateCommand{}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
	}
	os.Exit(exitStatus)
}
