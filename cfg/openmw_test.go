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

func TestRealOpenmwCFG(t *testing.T) {
	t.Skip()
	plugins, data, err := OpenMWPlugins("/home/ern/tes3/config/openmw.cfg")
	require.NoError(t, err)
	require.NotEmpty(t, plugins)
	require.NotEmpty(t, data)
}
