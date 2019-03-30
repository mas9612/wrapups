package wuserver

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

// WrapupsServer is the implementation of pb.WrapupsServer.
type WrapupsServer struct {
	client *elastic.Client
	index  string
}

// NewWrapupsServer creates and returns new WrapupsServer instance.
// This method also create index for Elasticsearch if necessary.
func NewWrapupsServer() (pb.WrapupsServer, error) {
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

	wuServer := &WrapupsServer{
		client: client,
		index:  defaultIndexName,
	}
	return wuServer, nil
}

// ListWrapups returns the list of wrapup document stored in Elasticsearch.
func (s *WrapupsServer) ListWrapups(context.Context, *pb.ListWrapupsRequest) (*pb.ListWrapupsResponse, error) {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize Elasticsearch client")
	}

	result, err := client.Search(s.index).Query(elastic.NewMatchAllQuery()).Do(context.Background())
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

// GetWrapup returns a wrapup document matched to request.
func (s *WrapupsServer) GetWrapup(context.Context, *pb.GetWrapupRequest) (*pb.Wrapup, error) {
	return nil, nil
}

// CreateWrapup creates new wrapup document and stores it in Elasticsearch.
func (s *WrapupsServer) CreateWrapup(ctx context.Context, req *pb.CreateWrapupRequest) (*pb.Wrapup, error) {
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
