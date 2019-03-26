package main

import (
	"context"
	"fmt"
	"os"

	pb "github.com/mas9612/wrapups/pkg/wrapups"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to gRPC server: %v\n", err)
		os.Exit(1)
	}

	client := pb.NewWrapupsClient(conn)

	req := &pb.ListWrapupsRequest{}
	res, err := client.ListWrapups(context.Background(), req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get response: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Count: %d\n", res.Count)
	for _, wrapup := range res.Wrapups {
		fmt.Printf("ID: %s\n", wrapup.Id)
		fmt.Printf("Title: %s\n", wrapup.Title)
		fmt.Printf("Wrapup: %s\n", wrapup.Wrapup)
		fmt.Printf("Comment: %s\n", wrapup.Comment)
		fmt.Printf("Note: %s\n\n", wrapup.Note)
	}
}
