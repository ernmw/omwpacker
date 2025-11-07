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

	t.Run("header", func(t *testing.T) {
		sub := records[0].GetSubrecord(HEDR)
		require.NotNil(t, sub)
		h := &Header{}
		h.Unmarshal(sub.Data)
		require.NotZero(t, h)
		require.Equal(t, float32(1.3), h.Version)
		require.Equal(t, uint32(1), h.NumRecords)
		require.Empty(t, h.Name)
		require.Empty(t, h.Description)

		h.Description = "changed"
		raw, err := h.Marshal()
		require.NoError(t, err)
		remarshaled := &Header{}
		remarshaled.Unmarshal(raw)
		require.Equal(t, "changed", remarshaled.Description)
	})
}
