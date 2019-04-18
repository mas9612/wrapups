package command

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	pb "github.com/mas9612/wrapups/pkg/wrapups"
	"google.golang.org/grpc"
)

// GetCommand implements get subcommand.
type GetCommand struct{}

// Help returns the long-form help text of get subcommand.
func (c *GetCommand) Help() string {
	helpText := `
Usage: wuclient get <id>
  Get wrapup document.

Options:
  -a, --address  Elasticsearch server address (default is localhost).
  -p, --port     Elasticsearch server port (default is 9200).
`
	return strings.TrimSpace(helpText)
}

type getOptions struct {
	Addr string `short:"a" long:"address" default:"localhost" description:"Wrapups server address. (default is localhost)"`
	Port int    `short:"p" long:"port" default:"10000" description:"Wrapups server port. (default is 10000)"`
	Args struct {
		ID string `description:"Wrapup document ID."`
	} `positional-args:"yes" required:"yes"`
}

// Run runs get subcommand and returns exit status.
func (c *GetCommand) Run(args []string) int {
	opts := getOptions{}
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.ParseArgs(args); err != nil {
		flagsErr := err.(*flags.Error)
		if flagsErr.Type == flags.ErrHelp {
			fmt.Printf("%s\n", flagsErr.Message)
			return 0
		}
		fmt.Fprintf(os.Stderr, "failed to parse command line flags: %s", err.Error())
		return 1
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", opts.Addr, opts.Port), grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to gRPC server: %v\n", err)
		return 1
	}
	defer conn.Close()
	client := pb.NewWrapupsClient(conn)

	req := &pb.GetWrapupRequest{
		Id: opts.Args.ID,
	}
	res, err := client.GetWrapup(context.Background(), req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get document: %v\n", err)
		return 1
	}

	printWrapup(res)

	return 0
}

// Synopsis returns one-line synopsis of get subcommamd.
func (c *GetCommand) Synopsis() string {
	return "Get wrapup document."
}
