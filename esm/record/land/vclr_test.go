package land

import (
	"bytes"
	"testing"

	"github.com/ernmw/omwpacker/esm"
)

func TestVCLRFieldUnmarshalMarshal(t *testing.T) {
	const size = vclrSize
	const depth = vclrDepth

	// Prepare flat input data: sequential bytes 0..(65*65*3-1)
	input := make([]byte, size*size*depth)
	for i := range input {
		input[i] = byte(i % 128) // keep in int8 range
	}

	// Create a 65x65 grid of ColorField pointers
	colors := make([][]*ColorField, size)
	for y := range size {
		colors[y] = make([]*ColorField, size)
		for x := range size {
			colors[y][x] = &ColorField{}
		}
	}

	// Wrap in VCLRField
	field := &VCLRField{Colors: colors}

	// Unmarshal from flat input
	sub := &esm.Subrecord{Tag: VCLR, Data: input}
	if err := field.Unmarshal(sub); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Check that the first few colors match input data
	for y := range 3 {
		for x := range 3 {
			vertex := field.Colors[y][x]
			offset := (y*size + x) * depth
			if vertex.GetR() != uint8(input[offset]) {
				t.Fatalf("vertex[%d][%d].X = %d, want %d", y, x, vertex.GetR(), input[offset])
			}
			if vertex.GetG() != uint8(input[offset+1]) {
				t.Fatalf("vertex[%d][%d].Y = %d, want %d", y, x, vertex.GetG(), input[offset+1])
			}
			if vertex.GetB() != uint8(input[offset+2]) {
				t.Fatalf("vertex[%d][%d].Z = %d, want %d", y, x, vertex.GetB(), input[offset+2])
			}
		}
	}

	// Marshal back to subrecord
	outSub, err := field.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	if outSub.Tag != VCLR {
		t.Fatalf("Marshal tag = %v, want %v", outSub.Tag, VCLR)
	}

	// Check that flattened bytes match input
	if !bytes.Equal(outSub.Data, input) {
		t.Fatal("Marshal output does not match original input data")
	}
}
