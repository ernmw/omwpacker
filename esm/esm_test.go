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
		require.NoError(t, sub.Unmarshal(h))
		require.NoError(t, err)
		require.NotNil(t, h)
		require.Equal(t, float32(1.3), h.Version)
		require.Equal(t, uint32(1), h.NumRecords)
		require.Empty(t, h.Name)
		require.Empty(t, h.Description)

	})
}

func TestHeader(t *testing.T) {
	h := &Header{
		Version:     1.3,
		Name:        "name",
		Description: "description",
		NumRecords:  1,
	}
	raw, err := h.Marshal()
	require.NoError(t, err)
	require.True(t, bytes.Contains(raw.Data, []byte("name")))
	require.True(t, bytes.Contains(raw.Data, []byte("description")))

	// in your test before raw.Unmarshal(h2)
	t.Logf("raw.Data len=%d", len(raw.Data))
	t.Logf("raw.Data[0:40] hex: % x", raw.Data[0:40])         // header start..name
	t.Logf("raw.Data[40:40+32] hex: % x", raw.Data[40:40+32]) // part of description start
	t.Logf("raw.Data[40:296] contains 'description'? %v", bytes.Contains(raw.Data[40:296], []byte("description")))

	h2 := &Header{}
	require.NoError(t, raw.Unmarshal(h2))
	require.Equal(t, "description", h2.Description)
	require.Equal(t, "name", h2.Name)
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
