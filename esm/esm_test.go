package esm_test

import (
	"bytes"
	"path"
	"slices"
	"testing"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/record/cell"
	"github.com/ernmw/omwpacker/esm/record/land"
	"github.com/ernmw/omwpacker/esm/record/lua"
	"github.com/ernmw/omwpacker/esm/record/tes3"
	"github.com/stretchr/testify/require"
)

func getSubrecord(r *esm.Record, tag esm.SubrecordTag) *esm.Subrecord {
	for _, s := range r.Subrecords {
		if s.Tag == tag {
			return s
		}
	}
	return nil
}

func TestLUAL(t *testing.T) {
	// read the test file
	inputFile := path.Join("testdata", "LUAL.omwaddon")
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
	reread, err := esm.ParsePluginData("lual.omwaddon", bytes.NewReader(written))
	require.NoError(t, err)
	require.Equal(t, records, reread)

	t.Run("header", func(t *testing.T) {
		sub := getSubrecord(records[0], tes3.HEDR)
		require.NotNil(t, sub)
		h := &tes3.HEDRdata{}
		require.NoError(t, sub.UnmarshalTo(h))
		require.NoError(t, err)
		require.NotNil(t, h)
		require.Equal(t, float32(1.3), h.Version)
		require.Equal(t, uint32(1), h.NumRecords)
		require.Empty(t, h.Name)
		require.Empty(t, h.Description)

	})
}

func TestCELL(t *testing.T) {
	// read the test file
	inputFile := path.Join("testdata", "CELL.omwaddon")
	records, err := esm.ParsePluginFile(inputFile)
	require.NoError(t, err)
	require.Len(t, records, 2)
	require.Equal(t, tes3.TES3, records[0].Tag)
	require.Equal(t, cell.CELL, records[1].Tag)

	// marshal
	var buff bytes.Buffer
	for _, rec := range records {
		require.NoError(t, rec.Write(&buff))
	}
	written := buff.Bytes()
	// unmarshal again
	reread, err := esm.ParsePluginData("cell.omwaddon", bytes.NewReader(written))
	require.NoError(t, err)
	require.Equal(t, records, reread)

	t.Run("header", func(t *testing.T) {
		sub := getSubrecord(records[0], tes3.HEDR)
		require.NotNil(t, sub)
		h := &tes3.HEDRdata{}
		require.NoError(t, sub.UnmarshalTo(h))
		require.NoError(t, err)
		require.NotNil(t, h)
		require.Equal(t, float32(1.3), h.Version)
		require.Equal(t, uint32(1), h.NumRecords)
		require.Empty(t, h.Name)
		require.Empty(t, h.Description)
	})

	t.Run("cell parsing", func(t *testing.T) {
		cellRec, err := cell.ParseCELL(records[1])
		require.NoError(t, err)
		require.NotNil(t, cellRec)
		require.Equal(t, "Balmora, Caius Cosades' House", cellRec.NAME.Value)
		ordered, err := cellRec.OrderedRecords()
		require.NoError(t, err)
		require.NotEmpty(t, ordered)
		require.Equal(t, records[1].Subrecords, ordered)
	})
}

func TestLAND(t *testing.T) {
	// read the test file
	inputFile := path.Join("testdata", "large.esp")
	records, err := esm.ParsePluginFile(inputFile)
	require.NoError(t, err)
	require.NotEmpty(t, records)
	require.Equal(t, tes3.TES3, records[0].Tag)

	// marshal
	var buff bytes.Buffer
	for _, rec := range records {
		require.NoError(t, rec.Write(&buff))
	}
	written := buff.Bytes()
	// unmarshal again
	reread, err := esm.ParsePluginData("large.esp", bytes.NewReader(written))
	require.NoError(t, err)
	require.Equal(t, records, reread)

	t.Run("header", func(t *testing.T) {
		sub := getSubrecord(records[0], tes3.HEDR)
		require.NotNil(t, sub)
		h := &tes3.HEDRdata{}
		require.NoError(t, sub.UnmarshalTo(h))
		require.NoError(t, err)
		require.NotNil(t, h)
		require.Equal(t, float32(1.3), h.Version)
	})

	landRecordIndex := slices.IndexFunc(records, func(rec *esm.Record) bool {
		return rec.Tag == land.LAND
	})
	require.Greater(t, landRecordIndex, -1)
	landRecord := records[landRecordIndex]
	require.NotNil(t, landRecord)

	t.Run("vhgt", func(t *testing.T) {
		vhgtRecordIndex := slices.IndexFunc(landRecord.Subrecords, func(rec *esm.Subrecord) bool {
			return rec.Tag == land.VHGT
		})
		require.Greater(t, vhgtRecordIndex, -1)
		vhgt := landRecord.Subrecords[vhgtRecordIndex]
		require.NotNil(t, vhgt)

		parsed := land.VHGTField{}
		require.NoError(t, parsed.Unmarshal(vhgt))

		heights := parsed.ComputeAbsoluteHeights()
		require.NotEmpty(t, heights)
	})

}
