package wuserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	debuglogger "github.com/mas9612/wrapups/pkg/logger"
	pb "github.com/mas9612/wrapups/pkg/wrapups"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultIndexName = "wrapups"
	typ              = "_doc"

	internalErrorMsg = "internal server error occured. please try again later."
)

// WrapupsServer is the implementation of pb.WrapupsServer.
type WrapupsServer struct {
	client *elastic.Client
	index  string
	logger *zap.Logger
}

type config struct {
	url   string
	port  int
	trace bool
}

// Option is wrapups server option.
type Option func(*config)

// SetURL sets Elasticsearch server url.
// Default is localhost.
func SetURL(url string) Option {
	return func(c *config) {
		c.url = url
	}
}

// SetPort sets Elasticsearch server port.
// Default is 9200.
func SetPort(port int) Option {
	return func(c *config) {
		c.port = port
	}
}

// SetTrace sets whether trace log is enabled.
// Default is false.
func SetTrace(trace bool) Option {
	return func(c *config) {
		c.trace = trace
	}
}

// NewWrapupsServer creates and returns new WrapupsServer instance.
// This method also create index for Elasticsearch if necessary.
func NewWrapupsServer(logger *zap.Logger, opts ...Option) (pb.WrapupsServer, error) {
	c := config{
		url:   "localhost",
		port:  9200,
		trace: false,
	}
	for _, o := range opts {
		o(&c)
	}

	wuServer := &WrapupsServer{
		logger: logger,
	}

	logger.Info("initializing Elasticsearch client")
	options := make([]elastic.ClientOptionFunc, 0, 10)
	options = append(options, elastic.SetSniff(false))
	if c.url != "localhost" || c.port != 9200 {
		options = append(options, elastic.SetURL(fmt.Sprintf("http://%s:%d", c.url, c.port)))
	}
	if c.trace {
		l := &debuglogger.Logger{
			Logger: logger,
		}
		options = append(options, elastic.SetTraceLog(l))
	}
	client, err := elastic.NewClient(options...)
	if err != nil {
		errMsg := "failed to initialize Elasticsearch client"
		logger.Error(errMsg, zap.Error(err))
		return nil, errors.Wrap(err, errMsg)
	}

	exists, err := client.IndexExists(defaultIndexName).Do(context.Background())
	if err != nil {
		errMsg := "failed to check whether index exists"
		logger.Error(errMsg, zap.Error(err))
		return nil, errors.Wrap(err, errMsg)
	}
	if !exists {
		logger.Info(fmt.Sprintf("index \"%s\" not found. creating", defaultIndexName))
		_, err := client.CreateIndex(defaultIndexName).Do(context.Background())
		if err != nil {
			errMsg := fmt.Sprintf("failed to create index \"%s\"", defaultIndexName)
			logger.Error(errMsg, zap.Error(err))
			return nil, errors.Wrapf(err, errMsg)
		}
	}

	wuServer.client = client
	wuServer.index = defaultIndexName
	logger.Info("server initialization finished")

	return wuServer, nil
}

// ListWrapups returns the list of wrapup document stored in Elasticsearch.
func (s *WrapupsServer) ListWrapups(ctx context.Context, req *pb.ListWrapupsRequest) (*pb.ListWrapupsResponse, error) {
	var query elastic.Query
	if req.Filter == "" {
		query = elastic.NewMatchAllQuery()
	} else {
		query = elastic.NewMatchQuery("wrapup", req.Filter)
	}
	result, err := s.client.Search(s.index).Query(query).Do(ctx)
	if err != nil {
		errMsg := "failed to get documents from Elasticsearch"
		s.logger.Error(errMsg, zap.Error(err))
		return nil, status.Error(codes.Internal, internalErrorMsg)
	}

	wrapups := make([]*pb.Wrapup, 0, result.TotalHits())
	for _, hit := range result.Hits.Hits {
		var wrapup pb.Wrapup
		if err := json.Unmarshal(*hit.Source, &wrapup); err != nil {
			errMsg := "failed to Unmarshal response to JSON"
			s.logger.Error(errMsg, zap.Error(err))
			return nil, status.Error(codes.Internal, internalErrorMsg)
		}
		wrapup.Id = hit.Id
		wrapups = append(wrapups, &wrapup)
	}

	return &pb.ListWrapupsResponse{
		Count:   int32(result.TotalHits()),
		Wrapups: wrapups,
	}, nil
}

// GetWrapup returns a wrapup document matched to request.
func (s *WrapupsServer) GetWrapup(ctx context.Context, req *pb.GetWrapupRequest) (*pb.Wrapup, error) {
	if req.Id == "" {
		errMsg := "Id is required"
		s.logger.Error(errMsg)
		return nil, status.Error(codes.InvalidArgument, errMsg)
	}

	result, err := s.client.Get().Index(s.index).Id(req.Id).Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			errMsg := fmt.Sprintf("ID %s not found", req.Id)
			return nil, status.Error(codes.NotFound, errMsg)
		}
		errMsg := "failed to get document from Elasticsearch"
		s.logger.Error(errMsg, zap.Error(err))
		return nil, status.Error(codes.Internal, internalErrorMsg)
	}

	doc := &pb.Wrapup{}
	if err := json.Unmarshal(*result.Source, doc); err != nil {
		errMsg := "failed to Unmarshal response to JSON"
		s.logger.Error(errMsg, zap.Error(err))
		return nil, status.Error(codes.Internal, internalErrorMsg)
	}
	doc.Id = result.Id
	return doc, nil
}

// CreateWrapup creates new wrapup document and stores it in Elasticsearch.
func (s *WrapupsServer) CreateWrapup(ctx context.Context, req *pb.CreateWrapupRequest) (*pb.Wrapup, error) {
	if req.Title == "" {
		errMsg := "Title is required"
		s.logger.Error(errMsg)
		return nil, status.Error(codes.InvalidArgument, errMsg)
	}
	// TODO: check if document which has requested title has already exist
	// if so, return codes.AlreadyExists

	r := struct {
		pb.CreateWrapupRequest
		CreateTime *timestamp.Timestamp `json:"create_time"`
	}{
		*req,
		ptypes.TimestampNow(),
	}
	res, err := s.client.Index().Index(s.index).Type(typ).BodyJson(r).Do(ctx)
	if err != nil {
		errMsg := "failed to create new document"
		s.logger.Error(errMsg, zap.Error(err))
		return nil, status.Error(codes.Internal, internalErrorMsg)
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
