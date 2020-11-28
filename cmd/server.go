package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-openapi/loads"
	"github.com/spf13/pflag"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/devchallenge/article-similarity/internal/article"
	"github.com/devchallenge/article-similarity/internal/http"
	"github.com/devchallenge/article-similarity/internal/http/restapi"
	"github.com/devchallenge/article-similarity/internal/http/restapi/operations"
	"github.com/devchallenge/article-similarity/internal/mongo"
	"github.com/devchallenge/article-similarity/internal/similarity"
)

const (
	defaultSimilarityThreshold = 0.95

	defaultStorageConnectTimeout = 10 * time.Second

	irregularVerbFilePath = "assets/irregular_verbs.csv"
)

type Config struct {
	SimilarityThreshold float64
	MongoHost           string
	MongoPort           int
	MongoDatabase       string
}

func (c *Config) InitFlags() {
	pflag.Float64Var(&c.SimilarityThreshold, "similarity_threshold", defaultSimilarityThreshold,
		"article similarity threshold in percents")
	pflag.StringVar(&c.MongoHost, "mongo_host", "localhost", "mongodb host")
	pflag.IntVar(&c.MongoPort, "mongo_port", 27017, "mongodb port")
	pflag.StringVar(&c.MongoDatabase, "mongo_database", "dev", "mongodb database name")
}

func ExecuteServer() error {
	config := &Config{}
	config.InitFlags()

	pflag.Parse()

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return fmt.Errorf("failed to embedded spec: %w", err)
	}

	api := operations.NewArticleSimilarityAPI(swaggerSpec)
	api.Logger = log.Printf
	rest := restapi.NewServer(api)

	defer func() {
		if serr := rest.Shutdown(); serr != nil {
			log.Printf("rest shutdown failed: %v", serr)
		}
	}()

	mongoURI := fmt.Sprintf("mongodb://%s:%d", config.MongoHost, config.MongoPort)
	log.Printf("mongoURI: %s", mongoURI)

	mc, err := mg.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("failed to create mongo: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageConnectTimeout)
	defer cancel()

	if err := mc.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	defer func() {
		if err := mc.Disconnect(context.Background()); err != nil {
			log.Printf("failed to disconnect mongo: %v", err)
		}
	}()

	st := mongo.New(mc, config.MongoDatabase)

	irregularVerb := similarity.IrregularVerb{}
	if err := irregularVerb.Load(irregularVerbFilePath); err != nil {
		log.Printf("failed to load irregular verbs from=%s: %v", irregularVerbFilePath, err)
	}

	art := article.New(similarity.NewSimilarity(config.SimilarityThreshold, irregularVerb), st)

	h := http.New(art)
	h.ConfigureHandlers(api)
	rest.ConfigureAPI()

	return rest.Serve()
}
