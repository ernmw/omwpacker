package record_test

import (
	"bytes"
	"testing"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/record"
	"github.com/stretchr/testify/require"
)

type anamTagger struct{}

func (t *anamTagger) Tag() esm.SubrecordTag { return "ANAM" }

type ANAMdata = record.ZstringSubrecord[*anamTagger]

func TestAbstract(t *testing.T) {
	a := ANAMdata{Value: "hello there"}
	sub, err := a.Marshal()
	require.NoError(t, err)
	require.NotNil(t, sub)
	require.True(t, bytes.Contains(sub.Data, []byte("hello there")))
	require.Equal(t, esm.SubrecordTag("ANAM"), sub.Tag)

	unmarsh := ANAMdata{}
	require.NoError(t, sub.Unmarshal(&unmarsh))
	require.Equal(t, a, unmarsh)
}
