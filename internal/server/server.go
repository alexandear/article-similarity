package server

import (
	"context"
	"fmt"
	"time"

	"github.com/go-openapi/loads"
	"github.com/hashicorp/go-multierror"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/devchallenge/article-similarity/internal/handler"
	"github.com/devchallenge/article-similarity/internal/mongo"
	"github.com/devchallenge/article-similarity/internal/server/article"
	"github.com/devchallenge/article-similarity/internal/similarity"
	"github.com/devchallenge/article-similarity/internal/swagger/restapi"
	"github.com/devchallenge/article-similarity/internal/swagger/restapi/operations"
)

const (
	defaultStorageConnectTimeout = 10 * time.Second

	irregularVerbFilePath = "assets/irregular_verbs.csv"
)

type Server struct {
	rest    *restapi.Server
	mongo   *mg.Client
	article *article.Article
}

func New(
	logger func(format string, v ...interface{}),
	mongoHost string, mongoPort int, mongoDatabase string,
	similarityThreshold float64,
) (*Server, error) {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to embedded spec: %w", err)
	}

	api := operations.NewArticleSimilarityAPI(swaggerSpec)
	api.Logger = logger
	rest := restapi.NewServer(api)

	mongoURI := fmt.Sprintf("mongodb://%s:%d", mongoHost, mongoPort)
	logger("mongoURI: %s", mongoURI)

	mc, err := mg.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageConnectTimeout)
	defer cancel()

	if err := mc.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	st := mongo.New(mc, mongoDatabase)

	irregularVerb := similarity.IrregularVerb{}
	if err := irregularVerb.Load(irregularVerbFilePath); err != nil {
		logger("failed to load irregular verbs from=%s: %v", irregularVerbFilePath, err)
	}

	art := article.New(logger, similarity.NewSimilarity(logger, similarityThreshold, irregularVerb), st)

	server := &Server{
		rest:    rest,
		mongo:   mc,
		article: art,
	}

	h := handler.New(art)
	h.ConfigureHandlers(api)
	rest.ConfigureAPI()

	return server, nil
}

func (s *Server) Serve() error {
	return s.rest.Serve()
}

func (s *Server) Close() error {
	var resErr error
	if err := s.rest.Shutdown(); err != nil {
		resErr = multierror.Append(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageConnectTimeout)
	defer cancel()

	if err := s.mongo.Disconnect(ctx); err != nil {
		resErr = multierror.Append(err)
	}

	return fmt.Errorf("failed to close: %w", resErr)
}
