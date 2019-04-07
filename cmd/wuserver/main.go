package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jessevdk/go-flags"
	pb "github.com/mas9612/wrapups/pkg/wrapups"
	"github.com/mas9612/wrapups/pkg/wuserver"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type options struct {
	Port        int    `short:"p" long:"port" description:"wrapups server port" default:"10000"`
	ElasticAddr string `long:"elastic-addr" description:"Elasticsearch server address (default: localhost)"`
	ElasticPort int    `long:"elastic-port" description:"Elasticsaerch server port (default: 9200)"`
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	opts := options{}
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		flagsErr := err.(*flags.Error)
		if flagsErr.Type == flags.ErrHelp {
			fmt.Printf("%s\n", flagsErr.Message)
			return
		}
		logger.Fatal("failed to parse command line flags", zap.Error(err))
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.Port))
	if err != nil {
		logger.Fatal("listen failed", zap.Error(err))
	}
	defer listener.Close()
	logger.Info(fmt.Sprintf("listening on :%d", opts.Port))

	wrapupsOpts := make([]wuserver.Option, 0, 5)
	if opts.ElasticAddr != "" {
		wrapupsOpts = append(wrapupsOpts, wuserver.SetURL(opts.ElasticAddr))
	}
	if opts.ElasticAddr != "" {
		wrapupsOpts = append(wrapupsOpts, wuserver.SetPort(opts.ElasticPort))
	}
	wuServer, err := wuserver.NewWrapupsServer(logger, wrapupsOpts...)
	if err != nil {
		logger.Fatal("server initialization failed", zap.Error(err))
	}
	grpcServer := grpc.NewServer()
	pb.RegisterWrapupsServer(grpcServer, wuServer)
	log.Fatal(grpcServer.Serve(listener))
}
