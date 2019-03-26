package main

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/mas9612/wrapups/pkg/wrapups"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

const (
	defaultIndexName = "wrapups"
	typ              = "_doc"
)

type wrapupsServer struct {
	client *elastic.Client
	index  string
}

func newWrapupsServer() (pb.WrapupsServer, error) {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize Elasticsearch client")
	}

	exists, err := client.IndexExists(defaultIndexName).Do(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to check whether index exists")
	}
	if !exists {
		_, err := client.CreateIndex(defaultIndexName).Do(context.Background())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create index \"%s\"", defaultIndexName)
		}
	}

	wuServer := &wrapupsServer{
		client: client,
		index:  defaultIndexName,
	}
	return wuServer, nil
}

func (s *wrapupsServer) ListWrapups(context.Context, *pb.ListWrapupsRequest) (*pb.ListWrapupsResponse, error) {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize Elasticsearch client")
	}

	result, err := client.Search("wrapups").Query(elastic.NewMatchAllQuery()).Do(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get documents from Elasticsearch")
	}

	wrapups := make([]*pb.Wrapup, 0, result.TotalHits())
	for _, hit := range result.Hits.Hits {
		var wrapup pb.Wrapup
		if err := json.Unmarshal(*hit.Source, &wrapup); err != nil {
			return nil, errors.Wrap(err, "failed to Unmarshal response to JSON")
		}
		wrapups = append(wrapups, &wrapup)
	}

	return &pb.ListWrapupsResponse{
		Count:   int32(result.TotalHits()),
		Wrapups: wrapups,
	}, nil
}

func (s *wrapupsServer) GetWrapup(context.Context, *pb.GetWrapupRequest) (*pb.Wrapup, error) {
	return nil, nil
}

func (s *wrapupsServer) CreateWrapup(ctx context.Context, req *pb.CreateWrapupRequest) (*pb.Wrapup, error) {
	if req.Title == "" {
		return nil, errors.New("Title is required")
	}

	r := struct {
		pb.CreateWrapupRequest
		CreateTime *timestamp.Timestamp
	}{
		*req,
		ptypes.TimestampNow(),
	}
	res, err := s.client.Index().Index(s.index).Type(typ).BodyJson(r).Do(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new document")
	}
	doc := &pb.Wrapup{
		Id:         res.Id,
		Title:      req.Title,
		Wrapup:     req.Wrapup,
		Comment:    req.Comment,
		Note:       req.Note,
		CreateTime: r.CreateTime,
	}
	return doc, nil
}
