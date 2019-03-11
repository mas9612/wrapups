package main

import (
	"context"

	pb "github.com/mas9612/wrapups/pkg/wrapups"
)

type wrapupsServer struct{}

func (s *wrapupsServer) ListWrapups(context.Context, *pb.ListWrapupsRequest) (*pb.ListWrapupsResponse, error) {
	return nil, nil
}

func (s *wrapupsServer) GetWrapup(context.Context, *pb.GetWrapupRequest) (*pb.Wrapup, error) {
	return nil, nil
}

func (s *wrapupsServer) CreateWrapup(context.Context, *pb.CreateWrapupRequest) (*pb.Wrapup, error) {
	return nil, nil
}
