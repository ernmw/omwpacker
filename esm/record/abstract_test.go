package record_test

import (
	"bytes"
	"testing"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/record"
	"github.com/stretchr/testify/require"
)

func TestAbstract(t *testing.T) {
	concType := record.Tag{Value: "CONC"}
	c := record.AbstractZSTRING[concType]

	c := &concreteZSTRING{AbstractZSTRING: "a value"}
	require.Equal(t, esm.SubrecordTag("CONC"), c.Tag())
	sub, err := c.Marshal()
	require.NoError(t, err)
	require.NotNil(t, sub)
	require.Equal(t, esm.SubrecordTag("CONC"), sub.Tag)
	require.True(t, bytes.Contains(sub.Data, []byte("a value")))
	unmarsh := concreteZSTRING{}
	require.NoError(t, sub.Unmarshal(&unmarsh))
	require.Equal(t, c, unmarsh)
}
