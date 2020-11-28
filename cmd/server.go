package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/pflag"

	"github.com/devchallenge/article-similarity/internal/server"
)

func ExecuteServer() error {
	config := &Config{}
	config.InitFlags()

	pflag.Parse()

	logger := log.Printf

	serv, err := server.New(logger, config.MongoHost, config.MongoPort, config.MongoDatabase,
		config.SimilarityThreshold)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	defer func() {
		if err := serv.Close(); err != nil {
			logger("server close failed: %v", err)
		}
	}()

	return serv.Serve()
}
