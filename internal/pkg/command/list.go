package command

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/jessevdk/go-flags"
	pb "github.com/mas9612/wrapups/pkg/wrapups"
	"google.golang.org/grpc"
)

// ListCommand implements list subcommand.
type ListCommand struct{}

// Help returns the long-form help text of list subcommand.
func (c *ListCommand) Help() string {
	helpText := `
Usage: wuclient list
  List wrapup documents.
`
	return strings.TrimSpace(helpText)
}

type options struct {
	Addr string `short:"a" long:"address" default:"localhost" description:"Wrapups server address. (default is localhost)"`
	Port int    `short:"p" long:"port" default:"10000" description:"Wrapups server port. (default is 10000)"`
}

// Run runs list subcommand and returns exit status.
func (c *ListCommand) Run(args []string) int {
	opts := options{}
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
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

	req := &pb.ListWrapupsRequest{}
	res, err := client.ListWrapups(context.Background(), req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get response: %v\n", err)
		return 1
	}

	fmt.Printf("Count: %d\n", res.Count)
	for _, wrapup := range res.Wrapups {
		fmt.Printf("ID: %s\n", wrapup.Id)
		fmt.Printf("Title: %s\n", wrapup.Title)
		fmt.Printf("Wrapup: %s\n", wrapup.Wrapup)
		fmt.Printf("Comment: %s\n", wrapup.Comment)
		fmt.Printf("Note: %s\n", wrapup.Note)
		t, err := ptypes.Timestamp(wrapup.CreateTime)
		if err != nil {
			fmt.Printf("CreateTime: <invalid>\n\n")
		} else {
			fmt.Printf("CreateTime: %s\n\n", t.String())
		}
	}

	return 0
}

// Synopsis returns one-line synopsis of list subcommamd.
func (c *ListCommand) Synopsis() string {
	return "List wrapup documents."
}
