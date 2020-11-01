package cmd

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitFlags(t *testing.T) {
	require.NoError(t, os.Setenv("PORT", "80"))
	require.NoError(t, os.Setenv("SIMILARITY_THRESHOLD", "0.92"))
	require.NoError(t, os.Setenv("MONGO_HOST", "mongo"))
	require.NoError(t, os.Setenv("MONGO_PORT", "27020"))
	os.Args = append(os.Args, "--port", "8050")
	config := &Config{}

	err := InitFlags(config)

	assert.NoError(t, err)
	assert.Equal(t, 0.92, config.SimilarityThreshold)
	assert.Equal(t, "mongo", config.MongoHost)
	assert.Equal(t, 27020, config.MongoPort)
	assertFlagEqual(t, "similarity_threshold", "0.92")
	assertFlagEqual(t, "port", "8050")
}

func assertFlagEqual(t *testing.T, name, expected string) {
	f := pflag.Lookup(name)
	require.NotNil(t, f)
	assert.Equal(t, expected, f.Value.String())
}
