package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/pflag"

	"github.com/devchallenge/article-similarity/internal/server"
	"github.com/devchallenge/article-similarity/internal/util"
	"github.com/devchallenge/article-similarity/internal/util/cmd"
)

func InitFlags(config *Config) error {
	config.InitFlags()

	pflag.Parse()

	if err := cmd.BindEnv(pflag.CommandLine); err != nil {
		return fmt.Errorf("failed to bind env: %w", err)
	}

	return nil
}

func ExecuteServer() error {
	config := &Config{}

	if err := InitFlags(config); err != nil {
		return fmt.Errorf("failed to init flags: %w", err)
	}

	logger := log.Printf

	serv, err := server.New(logger, config.MongoHost, config.MongoPort, config.MongoDatabase,
		config.SimilarityThreshold)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}
	defer util.Close(serv, logger)

	return serv.Serve()
}
