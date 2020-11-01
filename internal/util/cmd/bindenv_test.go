package cmd

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBindEnv(t *testing.T) {
	flags := pflag.NewFlagSet(t.Name(), pflag.PanicOnError)
	const flagName = "prefix.variable-suffix"
	var variable int
	flags.IntVar(&variable, flagName, 5, "")
	require.NoError(t, os.Setenv("PREFIX_VARIABLE_SUFFIX", "123"))
	pflag.Parse()

	assert.NoError(t, BindEnv(flags))
	actual := flags.Lookup(flagName)
	require.NotNil(t, actual)
	assert.Equal(t, "123", actual.Value.String())
}
