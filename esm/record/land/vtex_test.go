package land

import (
	"encoding/binary"
	"testing"

	"github.com/ernmw/omwpacker/esm"
)

func TestVTEXField_UnmarshalMarshal(t *testing.T) {
	const size = vtexSize
	const bytesPerEntry = 2
	totalBytes := size * size * bytesPerEntry

	// Build synthetic source data (each cell gets uint16 = its linear index)
	data := make([]byte, totalBytes)
	for i := range size * size {
		binary.LittleEndian.PutUint16(data[i*2:], uint16(i))
	}

	// --- Create and fill grid ---
	grid := make([][]uint16, size)
	for y := range size {
		grid[y] = make([]uint16, size)
		for x := range size {
			grid[y][x] = uint16(y*size + x)
		}
	}

	// --- Unmarshal ---
	v := &VTEXField{Vertices: make([][]uint16, size)}
	for y := range size {
		v.Vertices[y] = make([]uint16, size)
		for x := range size {
			v.Vertices[y][x] = 0
		}
	}

	sub := &esm.Subrecord{Tag: VTEX, Data: data}
	if err := v.Unmarshal(sub); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// --- Marshal ---
	out, err := v.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// --- Expected length check ---
	wantLen := totalBytes
	gotLen := len(out.Data)
	if gotLen != wantLen {
		t.Fatalf("expected data len %d, got %d", wantLen, gotLen)
	}

	// --- Byte comparison ---
	for i := range data {
		if out.Data[i] != data[i] {
			t.Fatalf("byte mismatch at %d: got=%02x want=%02x", i, out.Data[i], data[i])
		}
	}
}
