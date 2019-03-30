package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/mas9612/wrapups/pkg/wrapups"
	"github.com/mas9612/wrapups/pkg/wuserver"

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
	wuServer, err := wuserver.NewWrapupsServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "server initialization failed: %v\n", err)
		os.Exit(1)
	}
	pb.RegisterWrapupsServer(grpcServer, wuServer)
	log.Fatal(grpcServer.Serve(listener))
}
