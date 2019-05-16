package command

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mas9612/wrapups/pkg/auth"
	"github.com/mas9612/wrapups/pkg/config"
	pb "github.com/mas9612/wrapups/pkg/wrapups"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ListCommand implements list subcommand.
type ListCommand struct {
	Conf *config.Config
}

// Help returns the long-form help text of list subcommand.
func (c *ListCommand) Help() string {
	helpText := `
Usage: wuclient list
  List wrapup documents.
`
	return strings.TrimSpace(helpText)
}

// Run runs list subcommand and returns exit status.
func (c *ListCommand) Run(args []string) int {
	conn, err := grpc.Dial(c.Conf.WuserverURL, grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to gRPC server: %v\n", err)
		return 1
	}
	defer conn.Close()
	client := pb.NewWrapupsClient(conn)

	token, err := auth.Token()
	if err != nil {
		fmt.Fprintf(os.Stderr, "auth error: %s\n", err.Error())
		return 1
	}
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", fmt.Sprintf("bearer %s", token)))
	req := &pb.ListWrapupsRequest{}
	res, err := client.ListWrapups(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get response: %v\n", err)
		return 1
	}

	fmt.Printf("Count: %d\n", res.Count)
	for _, wrapup := range res.Wrapups {
		printWrapup(wrapup)
		fmt.Print("\n")
	}

	return 0
}

// Synopsis returns one-line synopsis of list subcommamd.
func (c *ListCommand) Synopsis() string {
	return "List wrapup documents."
}
