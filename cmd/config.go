package cmd

import (
	"github.com/spf13/pflag"
)

const (
	defaultSimilarityThreshold = 0.95
)

type Config struct {
	SimilarityThreshold float64
}

func (c *Config) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("config", pflag.PanicOnError)

	flags.Float64Var(&c.SimilarityThreshold, "similarity_threshold", defaultSimilarityThreshold,
		"article similarity threshold in percents")

	return flags
}
