package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/jessevdk/go-flags"
	"github.com/mas9612/wrapups/internal/pkg/command"
	"github.com/mas9612/wrapups/pkg/config"
	"github.com/mas9612/wrapups/pkg/version"
	"github.com/mitchellh/cli"
)

func main() {
	conf := config.ParseConfig()

	confFlag := config.Config{}
	parser := flags.NewParser(&confFlag, flags.PrintErrors|flags.PassDoubleDash|flags.IgnoreUnknown)
	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	// override if cli flags set
	if confFlag.AuthserverURL != "" {
		conf.AuthserverURL = confFlag.AuthserverURL
	}
	if confFlag.WuserverURL != "" {
		conf.WuserverURL = confFlag.WuserverURL
	}

	wrapHelpTextWithOptions := func(app string) cli.HelpFunc {
		fn := cli.BasicHelpFunc(app)
		return func(commands map[string]cli.CommandFactory) string {
			helpText := fn(commands)
			optionHelp := `
Options:
    --authserver-url    Authserver URL. Must include both address and port number. (default: "localhost:10000")
    --wuserver-url      Wrapups server URL. Must include both address and port number. (default: "localhost:10000")
`
			return helpText + strings.TrimRightFunc(optionHelp, unicode.IsSpace)
		}
	}

	app := "wuclient"
	c := cli.NewCLI(app, version.Version)
	c.Args = args
	c.HelpFunc = wrapHelpTextWithOptions(app)
	c.Commands = map[string]cli.CommandFactory{
		"list": func() (cli.Command, error) {
			return &command.ListCommand{Conf: conf}, nil
		},
		"get": func() (cli.Command, error) {
			return &command.GetCommand{Conf: conf}, nil
		},
		"create": func() (cli.Command, error) {
			return &command.CreateCommand{Conf: conf}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
	}
	os.Exit(exitStatus)
}
