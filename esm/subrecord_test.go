package esm

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeader(t *testing.T) {
	h := &HEDRdata{
		Version:     1.3,
		Name:        "name",
		Description: "description",
		NumRecords:  1,
	}
	raw, err := h.Marshal()
	require.NoError(t, err)
	require.True(t, bytes.Contains(raw.Data, []byte("name")))
	require.True(t, bytes.Contains(raw.Data, []byte("description")))

	h2 := &HEDRdata{}
	require.NoError(t, raw.Unmarshal(h2))
	require.Equal(t, "description", h2.Description)
	require.Equal(t, "name", h2.Name)
}

func TestLUAF(t *testing.T) {
	h := &LUAFdata{
		Target: "CREA",
	}
	raw, err := h.Marshal()
	require.NoError(t, err)
	require.True(t, bytes.Contains(raw.Data, []byte("CREA")))

	h2 := &LUAFdata{}
	require.NoError(t, raw.Unmarshal(h2))
	require.Equal(t, "CREA", h2.Target)

	t.Run("short tag", func(t *testing.T) {
		h := &LUAFdata{
			Target: "NPC",
		}
		raw, err := h.Marshal()
		require.NoError(t, err)
		require.True(t, bytes.Contains(raw.Data, []byte("NPC_")))

		h2 := &LUAFdata{}
		require.NoError(t, raw.Unmarshal(h2))
		require.Equal(t, "NPC_", h2.Target)
	})
}

func TestLUAS(t *testing.T) {
	h := &LUASdata{
		Path: "some/path/to/script.lua",
	}
	raw, err := h.Marshal()
	require.NoError(t, err)
	require.True(t, bytes.Contains(raw.Data, []byte("some/path/to/script.lua")))

	h2 := &LUASdata{}
	require.NoError(t, raw.Unmarshal(h2))
	require.Equal(t, "some/path/to/script.lua", h2.Path)
}
