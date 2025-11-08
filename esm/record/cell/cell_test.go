package cell

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/ernmw/omwpacker/esm"
)

func TestCellSubrecordMarshalUnmarshal(t *testing.T) {
	tests := []esm.Subrecord{
		{Tag: NAME, Data: []byte("Balmora")},
		{Tag: DATA, Data: func() []byte {
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, uint32(3))
			binary.Write(buf, binary.LittleEndian, int32(4))
			binary.Write(buf, binary.LittleEndian, int32(5))
			return buf.Bytes()
		}()},
		{Tag: RGNN, Data: []byte("Ascadian Isles")},
		{Tag: NAM5, Data: []byte{255, 128, 0}},
		{Tag: WHGT, Data: func() []byte {
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, float32(12.34))
			return buf.Bytes()
		}()},
		{Tag: AMBI, Data: func() []byte {
			buf := new(bytes.Buffer)
			// 3x3 = 9 color bytes
			buf.Write([]byte{10, 20, 30, 40, 50, 60, 70, 80, 90})
			// Add 3 padding bytes to total 12 before fog density
			buf.Write([]byte{0, 0, 0})
			// Fog density (4 bytes)
			binary.Write(buf, binary.LittleEndian, float32(0.5))
			return buf.Bytes()
		}()},
	}

	cell, err := ParseCELL(cloneSubs(tests))
	if err != nil {
		t.Fatalf("ParseCELL failed: %v", err)
	}

	// Collect marshaled output again
	marshalFuncs := []func() (*esm.Subrecord, error){
		cell.NAME.Marshal,
		cell.DATA.Marshal,
		cell.RGNN.Marshal,
		cell.NAM5.Marshal,
		cell.WHGT.Marshal,
		cell.AMBI.Marshal,
	}

	for i, fn := range marshalFuncs {
		if fn == nil {
			continue
		}
		sub, err := fn()
		if err != nil {
			t.Fatalf("Marshal #%d failed: %v", i, err)
		}

		orig := tests[i].Data
		got := sub.Data
		if !bytes.Equal(orig, got) {
			t.Errorf("%s: data mismatch\norig=%v\ngot =%v", tests[i].Tag, orig, got)
		}
	}
}

func cloneSubs(in []esm.Subrecord) []*esm.Subrecord {
	out := make([]*esm.Subrecord, len(in))
	for i := range in {
		cpy := in[i]
		out[i] = &cpy
	}
	return out
}
