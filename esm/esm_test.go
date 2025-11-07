package esm

import (
	"bytes"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	// read the test file
	inputFile := path.Join("testdata", "lua_conf_test.omwaddon")
	records, err := ParsePluginFile(inputFile)
	require.NoError(t, err)
	require.Len(t, records, 2)
	require.Equal(t, TES3, records[0].Tag)
	require.Equal(t, LUAL, records[1].Tag)

	// marshal
	var buff bytes.Buffer
	for _, rec := range records {
		require.NoError(t, rec.Write(&buff))
	}
	written := buff.Bytes()
	// unmarshal again
	reread, err := ParsePluginData("lua_conf_test.omwaddon", bytes.NewReader(written))
	require.NoError(t, err)
	require.Equal(t, records, reread)
}
