package land

import (
	"bytes"
	"testing"

	"github.com/ernmw/omwpacker/esm"
)

func TestWNAMField_UnmarshalMarshal(t *testing.T) {
	const size = wnamSize // 9
	const bytesPerEntry = 1

	totalBytes := size * size * bytesPerEntry

	// --- Build synthetic source data ---
	data := make([]byte, totalBytes)
	for i := range size * size {
		data[i] = uint8(i)
	}

	sub := &esm.Subrecord{Tag: WNAM, Data: data}

	// --- Allocate grid for Unmarshal ---
	grid := make([][]*ByteField, size)
	for y := range size {
		grid[y] = make([]*ByteField, size)
		for x := range size {
			grid[y][x] = new(ByteField)
		}
	}

	field := &WNAMField{Heights: grid}

	// --- Unmarshal ---
	if err := field.Unmarshal(sub); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// --- Verify a few sample values ---
	for y := range 3 {
		for x := range 3 {
			got := *field.Heights[y][x]
			want := ByteField(y*size + x)
			if got != want {
				t.Fatalf("Heights[%d][%d] = %d, want %d", y, x, got, want)
			}
		}
	}

	// --- Marshal ---
	outSub, err := field.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if outSub.Tag != WNAM {
		t.Fatalf("Tag = %v, want %v", outSub.Tag, WNAM)
	}

	if len(outSub.Data) != len(data) {
		t.Fatalf("output len = %d, want %d", len(outSub.Data), len(data))
	}

	if !bytes.Equal(outSub.Data, data) {
		t.Fatal("marshaled data mismatch")
	}
}
