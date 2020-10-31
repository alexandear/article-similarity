package cmd

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/devchallenge/article-similarity/internal/restapi"
	"github.com/devchallenge/article-similarity/internal/util"
	"github.com/devchallenge/article-similarity/internal/util/cmd"
)

func InitFlags(config *Config) {
	if err := cmd.BindEnv(pflag.CommandLine); err != nil {
		panic(err)
	}

	config.InitFlags()

	pflag.Parse()
}

func ExecuteServer() error {
	config := &Config{}

	InitFlags(config)

	serv, err := restapi.NewArticleServer(log.Printf, config.MongoHost, config.SimilarityThreshold)
	if err != nil {
		return errors.WithStack(err)
	}
	defer util.Close(serv)

	return serv.Serve()
}
