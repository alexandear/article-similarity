package cmd

import (
	"github.com/spf13/pflag"
)

const (
	defaultSimilarityThreshold = 0.95
)

type Config struct {
	SimilarityThreshold float64
	MongoHost           string
}

func (c *Config) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("config", pflag.PanicOnError)

	flags.Float64Var(&c.SimilarityThreshold, "similarity_threshold", defaultSimilarityThreshold,
		"article similarity threshold in percents")
	flags.StringVar(&c.MongoHost, "mongo_host", "mongo", "mongodb host")

	return flags
}
