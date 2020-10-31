package cmd

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitFlags(t *testing.T) {
	require.NoError(t, os.Setenv("PORT", "8080"))
	os.Args = append(os.Args, "--similarity_threshold", "0.90", "--mongo_host", "mongo")

	config := &Config{}
	InitFlags(config)

	assert.Equal(t, 0.90, config.SimilarityThreshold)
	assert.Equal(t, "mongo", config.MongoHost)
	port := pflag.Lookup("port")
	require.NotNil(t, port)
	assert.Equal(t, "8080", port.Value.String())
}
