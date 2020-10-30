package restapi

import (
	"github.com/go-openapi/loads"
	"github.com/kelseyhightower/memkv"
	"github.com/pkg/errors"

	"github.com/devchallenge/article-similarity/internal/handler"
	"github.com/devchallenge/article-similarity/internal/restapi/operations"
	"github.com/devchallenge/article-similarity/internal/similarity"
)

type ArticleServer struct {
	rest *Server
}

func NewArticleServer(logger func(string, ...interface{}), similarityThreshold float64) (*ArticleServer, error) {
	swaggerSpec, err := loads.Embedded(SwaggerJSON, FlatSwaggerJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to embedded spec")
	}

	api := operations.NewArticleSimilarityAPI(swaggerSpec)
	rest := NewServer(api)
	server := &ArticleServer{
		rest: rest,
	}

	store := memkv.New()
	sim := similarity.NewSimilarity(logger, similarityThreshold)

	h := handler.New(&store, sim)
	h.ConfigureHandlers(api)
	rest.ConfigureAPI()
	rest.api.Logger = logger

	return server, nil
}

func (s *ArticleServer) Serve() error {
	return s.rest.Serve()
}

func (s *ArticleServer) Close() error {
	return s.rest.Shutdown()
}
