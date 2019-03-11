package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/mas9612/wrapups/pkg/wrapups"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":10000")
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen failed: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterWrapupsServer(grpcServer, &wrapupsServer{})
	log.Fatal(grpcServer.Serve(listener))
}
