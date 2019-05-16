package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mas9612/wrapups/internal/pkg/command"
	"github.com/mas9612/wrapups/pkg/config"
	"github.com/mas9612/wrapups/pkg/version"
	"github.com/mitchellh/cli"
)

func main() {
	conf := config.Config{}
	parser := flags.NewParser(&conf, flags.PrintErrors|flags.PassDoubleDash|flags.IgnoreUnknown)
	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	// TODO: add top level options into help text
	c := cli.NewCLI("wuclient", version.Version)
	c.Args = args
	c.Commands = map[string]cli.CommandFactory{
		"list": func() (cli.Command, error) {
			return &command.ListCommand{Conf: &conf}, nil
		},
		"get": func() (cli.Command, error) {
			return &command.GetCommand{Conf: &conf}, nil
		},
		"create": func() (cli.Command, error) {
			return &command.CreateCommand{Conf: &conf}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
	}
	os.Exit(exitStatus)
}
