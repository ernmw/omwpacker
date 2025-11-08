package esm_test

import (
	"bytes"
	"path"
	"testing"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/record/lua"
	"github.com/ernmw/omwpacker/esm/record/tes3"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	// read the test file
	inputFile := path.Join("testdata", "lua_conf_test.omwaddon")
	records, err := esm.ParsePluginFile(inputFile)
	require.NoError(t, err)
	require.Len(t, records, 2)
	require.Equal(t, tes3.TES3, records[0].Tag)
	require.Equal(t, lua.LUAL, records[1].Tag)

	// marshal
	var buff bytes.Buffer
	for _, rec := range records {
		require.NoError(t, rec.Write(&buff))
	}
	written := buff.Bytes()
	// unmarshal again
	reread, err := esm.ParsePluginData("lua_conf_test.omwaddon", bytes.NewReader(written))
	require.NoError(t, err)
	require.Equal(t, records, reread)

	t.Run("header", func(t *testing.T) {
		sub := records[0].GetSubrecord(tes3.HEDR)
		require.NotNil(t, sub)
		h := &tes3.HEDRdata{}
		require.NoError(t, sub.Unmarshal(h))
		require.NoError(t, err)
		require.NotNil(t, h)
		require.Equal(t, float32(1.3), h.Version)
		require.Equal(t, uint32(1), h.NumRecords)
		require.Empty(t, h.Name)
		require.Empty(t, h.Description)

	})
}
