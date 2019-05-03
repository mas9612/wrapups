package command

import (
	"fmt"
	"strings"

	"github.com/mas9612/wrapups/pkg/version"
)

// VersionCommand implements version subcommand.
type VersionCommand struct{}

// Help returns the long-form help text of version subcommand.
func (c *VersionCommand) Help() string {
	helpText := `
Usage: wuclient version
  Print wrapups version.
`
	return strings.TrimSpace(helpText)
}

// Run runs version subcommand and returns exit status.
func (c *VersionCommand) Run(args []string) int {
	fmt.Println(version.Version)
	return 0
}

// Synopsis returns one-line synopsis of version subcommamd.
func (c *VersionCommand) Synopsis() string {
	return "Print wrapups version."
}
