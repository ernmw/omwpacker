package util

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrings(t *testing.T) {
	t.Run("fits", func(t *testing.T) {
		var out bytes.Buffer
		require.NoError(t, WritePaddedString(&out, []byte("howdy"), 20))
		got := ReadPaddedString(out.Bytes())
		require.Equal(t, "howdy", got)
	})
	t.Run("max size", func(t *testing.T) {
		var out bytes.Buffer
		require.NoError(t, WritePaddedString(&out, []byte("123"), 3))
		got := ReadPaddedString(out.Bytes())
		require.Equal(t, "123", got)
	})
	t.Run("too big", func(t *testing.T) {
		var out bytes.Buffer
		require.Error(t, WritePaddedString(&out, []byte("1234"), 3))
	})
}
