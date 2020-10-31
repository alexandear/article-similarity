package cmd

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/devchallenge/article-similarity/internal/restapi"
	"github.com/devchallenge/article-similarity/internal/util"
)

func InitFlags(config *Config) {
	config.InitFlags()

	pflag.Parse()
}

func ExecuteServer() error {
	config := &Config{}

	InitFlags(config)

	serv, err := restapi.NewArticleServer(log.Printf, config.SimilarityThreshold)
	if err != nil {
		return errors.WithStack(err)
	}
	defer util.Close(serv)

	return serv.Serve()
}
