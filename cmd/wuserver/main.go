package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/dgrijalva/jwt-go"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/jessevdk/go-flags"
	"github.com/mas9612/authserver/pkg/server"
	"github.com/mas9612/wrapups/pkg/version"
	pb "github.com/mas9612/wrapups/pkg/wrapups"
	"github.com/mas9612/wrapups/pkg/wuserver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type options struct {
	Port        int    `short:"p" long:"port" description:"wrapups server port" default:"10000"`
	ElasticAddr string `long:"elastic-addr" default:"localhost" description:"Elasticsearch server address"`
	ElasticPort int    `long:"elastic-port" default:"9200" description:"Elasticsaerch server port"`
	TraceLog    bool   `long:"trace" description:"Enable trace log."`
	Version     bool   `short:"v" long:"version" description:"Print wrapups version"`
}

func main() {
	l, _ := zap.NewProduction()
	defer l.Sync()

	opts := options{}
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		flagsErr := err.(*flags.Error)
		if flagsErr.Type == flags.ErrHelp {
			fmt.Printf("%s\n", flagsErr.Message)
			return
		}
		l.Fatal("failed to parse command line flags", zap.Error(err))
	}

	if opts.Version {
		fmt.Println(version.Version)
		return
	}

	var logger *zap.Logger
	if opts.TraceLog {
		config := zap.NewProductionConfig()
		config.Level.SetLevel(zapcore.DebugLevel)
		var err error
		if logger, err = config.Build(); err != nil {
			l.Fatal("failed to build logger", zap.Error(err))
		}
	} else {
		logger = l
	}
	defer logger.Sync()

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
	if opts.ElasticPort != 9200 {
		wrapupsOpts = append(wrapupsOpts, wuserver.SetPort(opts.ElasticPort))
	}
	if opts.TraceLog {
		wrapupsOpts = append(wrapupsOpts, wuserver.SetTrace(opts.TraceLog))
	}
	wuServer, err := wuserver.NewWrapupsServer(logger, wrapupsOpts...)
	if err != nil {
		logger.Fatal("server initialization failed", zap.Error(err))
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(authFunc),
		)),
	)
	pb.RegisterWrapupsServer(grpcServer, wuServer)
	log.Fatal(grpcServer.Serve(listener))
}

func authFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	claim := server.AuthClaim{}
	_, err = jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("requested signing method is not supported")
		}

		b, err := ioutil.ReadFile("./authserver.pub")
		if err != nil {
			return nil, err
		}
		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
		if err != nil {
			return nil, err
		}
		return verifyKey, nil
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("failed to verify token: %s", err.Error()))
	}

	return context.WithValue(ctx, "user", claim.User), nil
}
