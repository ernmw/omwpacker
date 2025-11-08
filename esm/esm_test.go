package esm

import (
	"bytes"
	"path"
	"testing"

	"github.com/ernmw/omwpacker/esm/tags"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	// read the test file
	inputFile := path.Join("testdata", "lua_conf_test.omwaddon")
	records, err := ParsePluginFile(inputFile)
	require.NoError(t, err)
	require.Len(t, records, 2)
	require.Equal(t, tags.TES3, records[0].Tag)
	require.Equal(t, tags.LUAL, records[1].Tag)

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
		sub := records[0].GetSubrecord(tags.HEDR)
		require.NotNil(t, sub)
		h := &HEDRdata{}
		require.NoError(t, sub.Unmarshal(h))
		require.NoError(t, err)
		require.NotNil(t, h)
		require.Equal(t, float32(1.3), h.Version)
		require.Equal(t, uint32(1), h.NumRecords)
		require.Empty(t, h.Name)
		require.Empty(t, h.Description)

	})
}

func TestStrings(t *testing.T) {
	t.Run("fits", func(t *testing.T) {
		var out bytes.Buffer
		require.NoError(t, writePaddedString(&out, []byte("howdy"), 20))
		got := readPaddedString(out.Bytes())
		require.Equal(t, "howdy", got)
	})
	t.Run("max size", func(t *testing.T) {
		var out bytes.Buffer
		require.NoError(t, writePaddedString(&out, []byte("123"), 3))
		got := readPaddedString(out.Bytes())
		require.Equal(t, "123", got)
	})
	t.Run("too big", func(t *testing.T) {
		var out bytes.Buffer
		require.Error(t, writePaddedString(&out, []byte("1234"), 3))
	})
}
