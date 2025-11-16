package cfg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const cfgPath = "/home/ern/tes3/config/openmw.cfg"
const bsaPath = "/home/ern/tes3/Morrowind/Data Files/Bloodmoon.bsa"

func TestOpenmwCFG(t *testing.T) {
	env, err := Load("./testdata")
	require.NoError(t, err)
	require.NotEmpty(t, env.Plugins)
	require.NotEmpty(t, env.Data)
}

func TestRealOpenmwCFG(t *testing.T) {
	if _, err := os.Stat(cfgPath); err != nil {
		t.Skip("CFG not present")
	}
	env, err := Load(cfgPath)
	require.NoError(t, err)
	require.NotEmpty(t, env.Plugins)
	require.NotEmpty(t, env.Data)

	require.NotEmpty(t, env.BSA)
	require.Contains(t, env.BSA, bsaPath)

	raw, err := env.ReadFile("textures/tx_bm_dirt_snow_01.dds")
	require.NoError(t, err)
	require.NotEmpty(t, raw)
	require.NoError(t, os.WriteFile("testdata/tx_bm_dirt_snow_01.dds", raw, 0666))
}
