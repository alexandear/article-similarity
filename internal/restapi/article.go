package restapi

import (
	"context"
	"fmt"
	"time"

	"github.com/go-openapi/loads"
	"github.com/hashicorp/go-multierror"
	"github.com/kelseyhightower/memkv"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/devchallenge/article-similarity/internal/handler"
	"github.com/devchallenge/article-similarity/internal/restapi/operations"
	"github.com/devchallenge/article-similarity/internal/similarity"
)

const (
	defaultStorageConnectTimeout = 10 * time.Second
)

type ArticleServer struct {
	rest  *Server
	mongo *mongo.Client
}

func NewArticleServer(logger func(string, ...interface{}), mongoHost string, similarityThreshold float64,
) (*ArticleServer, error) {
	swaggerSpec, err := loads.Embedded(SwaggerJSON, FlatSwaggerJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to embedded spec")
	}

	api := operations.NewArticleSimilarityAPI(swaggerSpec)
	rest := NewServer(api)

	mongoURI := fmt.Sprintf("mongodb://%s:27017", mongoHost)
	logger("mongoURI: %s", mongoURI)

	mc, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create mongo")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageConnectTimeout)
	defer cancel()

	if err := mc.Connect(ctx); err != nil {
		return nil, errors.WithStack(err)
	}

	server := &ArticleServer{
		rest:  rest,
		mongo: mc,
	}

	store := memkv.New()
	sim := similarity.NewSimilarity(logger, similarityThreshold)

	h := handler.New(mc, &store, sim)
	h.ConfigureHandlers(api)
	rest.ConfigureAPI()
	rest.api.Logger = logger

	return server, nil
}

func (s *ArticleServer) Serve() error {
	return s.rest.Serve()
}

func (s *ArticleServer) Close() error {
	var resErr error
	if err := s.rest.Shutdown(); err != nil {
		resErr = multierror.Append(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageConnectTimeout)
	defer cancel()

	if err := s.mongo.Disconnect(ctx); err != nil {
		resErr = multierror.Append(err)
	}

	return errors.WithStack(resErr)
}
