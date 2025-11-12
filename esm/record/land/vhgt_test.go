package land

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/ernmw/omwpacker/esm"
)

func TestVHGTField_UnmarshalMarshal(t *testing.T) {
	const size = vhgtSize

	// --- Construct synthetic input buffer ---

	offsetVal := float32(123.456)
	gridBytes := make([]byte, size*size)
	for i := range gridBytes {
		gridBytes[i] = byte(i % 127) // int8-compatible test data
	}
	junk := []byte{0xAA, 0xBB, 0xCC}

	buf := make([]byte, 4+len(gridBytes)+len(junk))
	binary.LittleEndian.PutUint32(buf[0:4], math.Float32bits(offsetVal))
	copy(buf[4:], gridBytes)
	copy(buf[4+len(gridBytes):], junk)

	sub := &esm.Subrecord{Tag: VHGT, Data: buf}

	// --- Allocate target struct and backing grid ---
	heights := make([][]uint8, size)
	for y := range size {
		heights[y] = make([]uint8, size)
		for x := range size {
			heights[y][x] = 0
		}
	}

	field := &VHGTField{Heights: heights}

	// --- Unmarshal ---
	if err := field.Unmarshal(sub); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Check offset decoded correctly
	if field.Offset != offsetVal {
		t.Fatalf("Offset = %v, want %v", field.Offset, offsetVal)
	}

	// Verify some sample height values
	for y := range 3 {
		for x := range 3 {
			got := byte(field.Heights[y][x])
			want := gridBytes[y*size+x]
			if got != want {
				t.Fatalf("Heights[%d][%d] = %d, want %d", y, x, got, want)
			}
		}
	}

	// --- Marshal back ---
	outSub, err := field.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if outSub.Tag != VHGT {
		t.Fatalf("Tag = %v, want %v", outSub.Tag, VHGT)
	}

	// --- Verify round-trip binary data ---
	if len(outSub.Data) != len(buf) {
		t.Fatalf("output len=%d, want %d", len(outSub.Data), len(buf))
	}

	if !bytes.Equal(outSub.Data[:len(buf)-3], buf[:len(buf)-3]) {
		t.Fatal("marshaled data mismatch (excluding junk)")
	}
}

func TestComputeAbsoluteHeights(t *testing.T) {
	// Make a synthetic VHGTField where all deltas are 1
	rows := make([][]uint8, vhgtSize)
	for y := range vhgtSize {
		row := make([]uint8, vhgtSize)
		for x := range vhgtSize {
			row[x] = 1
		}
		rows[y] = row
	}

	v := &VHGTField{
		Offset:  0.0,
		Heights: rows,
	}

	hmap := v.ComputeAbsoluteHeights()
	if hmap[0][0] != 8.0 {
		t.Fatalf("expected 8.0, got %v", hmap[0][0])
	}
	if hmap[vhgtSize-1][vhgtSize-1] <= 0 {
		t.Fatalf("expected positive height at end, got %v", hmap[vhgtSize-1][vhgtSize-1])
	}
}
