package tes3

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
	require.NoError(t, raw.UnmarshalTo(h2))
	require.Equal(t, "description", h2.Description)
	require.Equal(t, "name", h2.Name)
}
