package cfg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenmwCFG(t *testing.T) {
	plugins, data, err := OpenMWPlugins("./testdata")
	require.NoError(t, err)
	require.NotEmpty(t, plugins)
	require.NotEmpty(t, data)
}
