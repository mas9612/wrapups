package command

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	pb "github.com/mas9612/wrapups/pkg/wrapups"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

// CreateCommand implements create subcommand.
type CreateCommand struct{}

// Help returns the long-form help text of create subcommand.
func (c *CreateCommand) Help() string {
	helpText := `
Usage: wuclient create -f <filename>
  Create new wrapup document.

Options:
  -a, --address  Elasticsearch server address (default is localhost).
  -p, --port     Elasticsearch server port (default is 9200).
  -f, --file  Input filename. Required.
`
	return strings.TrimSpace(helpText)
}

type createOptions struct {
	Addr     string `short:"a" long:"address" default:"localhost" description:"Wrapups server address. (default is localhost)"`
	Port     int    `short:"p" long:"port" default:"10000" description:"Wrapups server port. (default is 10000)"`
	Filename string `short:"f" long:"file" required:"yes" description:"Input filename. Required."`
}

type yamlData struct {
	Title    string `yaml:"title,omitempty"`
	Wrapup   string `yaml:"wrapup,omitempty"`
	Comments string `yaml:"comments,omitempty"`
	Notes    string `yaml:"notes,omitempty"`
}

// Run runs create subcommand and returns exit status.
func (c *CreateCommand) Run(args []string) int {
	opts := createOptions{}
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

	b, err := ioutil.ReadFile(opts.Filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file: %v\n", err)
		return 1
	}
	var data yamlData
	if err := yaml.Unmarshal(b, &data); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse YAML file: %v\n", err)
		return 1
	}

	req := &pb.CreateWrapupRequest{
		Title:   data.Title,
		Wrapup:  data.Wrapup,
		Comment: data.Comments,
		Note:    data.Notes,
	}
	res, err := client.CreateWrapup(context.Background(), req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create document: %v\n", err)
		return 1
	}
	fmt.Printf("ID \"%s\" created\n", res.Id)

	return 0
}

// Synopsis returns one-line synopsis of create subcommamd.
func (c *CreateCommand) Synopsis() string {
	return "Create new wrapup document."
}
