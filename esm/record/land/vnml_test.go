package land

import (
	"bytes"
	"testing"

	"github.com/ernmw/omwpacker/esm"
)

func TestVNMLFieldUnmarshalMarshal(t *testing.T) {
	const size = vnmlSize
	const depth = vnmlDepth

	// Prepare flat input data: sequential bytes 0..(65*65*3-1)
	input := make([]byte, size*size*depth)
	for i := range input {
		input[i] = byte(i % 128) // keep in int8 range
	}

	// Create a 65x65 grid of VertexField pointers
	vertices := make([][]*VertexField, size)
	for y := range size {
		vertices[y] = make([]*VertexField, size)
		for x := range size {
			vertices[y][x] = &VertexField{}
		}
	}

	// Wrap in VNMLField
	field := &VNMLField{Vertices: vertices}

	// Unmarshal from flat input
	sub := &esm.Subrecord{Tag: VNML, Data: input}
	if err := field.Unmarshal(sub); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Check that the first few vertices match input data
	for y := range 3 {
		for x := range 3 {
			vertex := field.Vertices[y][x]
			offset := (y*size + x) * depth
			if vertex.GetX() != int8(input[offset]) {
				t.Fatalf("vertex[%d][%d].X = %d, want %d", y, x, vertex.GetX(), input[offset])
			}
			if vertex.GetY() != int8(input[offset+1]) {
				t.Fatalf("vertex[%d][%d].Y = %d, want %d", y, x, vertex.GetY(), input[offset+1])
			}
			if vertex.GetZ() != int8(input[offset+2]) {
				t.Fatalf("vertex[%d][%d].Z = %d, want %d", y, x, vertex.GetZ(), input[offset+2])
			}
		}
	}

	// Marshal back to subrecord
	outSub, err := field.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	if outSub.Tag != VNML {
		t.Fatalf("Marshal tag = %v, want %v", outSub.Tag, VNML)
	}

	// Check that flattened bytes match input
	if !bytes.Equal(outSub.Data, input) {
		t.Fatal("Marshal output does not match original input data")
	}
}
