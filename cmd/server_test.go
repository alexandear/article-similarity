package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitFlags(t *testing.T) {
	os.Args = append(os.Args, "--similarity_threshold", "0.90", "--port", "80")

	config := &Config{}
	InitFlags(config)

	assert.Equal(t, 0.90, config.SimilarityThreshold)
}
