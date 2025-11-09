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

type stringData = record.ZstringSubrecord[*anamTagger]

type floatData = record.Float32Subrecord[*anamTagger]

type intData = record.Uint32Subrecord[*anamTagger]

type byteData = record.Uint8Subrecord[*anamTagger]

func TestZstring(t *testing.T) {
	a := stringData{Value: "hello there"}
	sub, err := a.Marshal()
	require.NoError(t, err)
	require.NotNil(t, sub)
	require.True(t, bytes.Contains(sub.Data, []byte("hello there")))
	require.Equal(t, esm.SubrecordTag("ANAM"), sub.Tag)

	unmarsh := stringData{}
	require.NoError(t, sub.Unmarshal(&unmarsh))
	require.Equal(t, a, unmarsh)
}

func TestFloat32(t *testing.T) {
	a := floatData{Value: 5.0}
	sub, err := a.Marshal()
	require.NoError(t, err)
	require.NotNil(t, sub)
	require.Equal(t, esm.SubrecordTag("ANAM"), sub.Tag)

	unmarsh := floatData{}
	require.NoError(t, sub.Unmarshal(&unmarsh))
	require.Equal(t, a, unmarsh)
}

func TestUint32(t *testing.T) {
	a := intData{Value: 5}
	sub, err := a.Marshal()
	require.NoError(t, err)
	require.NotNil(t, sub)
	require.Equal(t, esm.SubrecordTag("ANAM"), sub.Tag)

	unmarsh := intData{}
	require.NoError(t, sub.Unmarshal(&unmarsh))
	require.Equal(t, a, unmarsh)
}

func TestUint8(t *testing.T) {
	a := byteData{Value: 5}
	sub, err := a.Marshal()
	require.NoError(t, err)
	require.NotNil(t, sub)
	require.Equal(t, esm.SubrecordTag("ANAM"), sub.Tag)

	unmarsh := byteData{}
	require.NoError(t, sub.Unmarshal(&unmarsh))
	require.Equal(t, a, unmarsh)
}
