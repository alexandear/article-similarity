package restapi

import (
	"github.com/go-openapi/loads"
	"github.com/kelseyhightower/memkv"
	"github.com/pkg/errors"

	"github.com/devchallenge/article-similarity/internal/handler"
	"github.com/devchallenge/article-similarity/internal/restapi/operations"
)

type ArticleServer struct {
	rest *Server
}

func NewArticleServer() (*ArticleServer, error) {
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

	h := handler.New(&store)
	h.ConfigureHandlers(api)
	rest.ConfigureAPI()

	return server, nil
}

func (s *ArticleServer) ConfigureLogger(logger func(string, ...interface{})) {
	s.rest.api.Logger = logger
}

func (s *ArticleServer) Serve() error {
	return s.rest.Serve()
}

func (s *ArticleServer) Close() error {
	return s.rest.Shutdown()
}
