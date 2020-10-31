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

func (c *Config) InitFlags() {
	pflag.Float64Var(&c.SimilarityThreshold, "similarity_threshold", defaultSimilarityThreshold,
		"article similarity threshold in percents")
	pflag.StringVar(&c.MongoHost, "mongo_host", "localhost", "mongodb host")
}
