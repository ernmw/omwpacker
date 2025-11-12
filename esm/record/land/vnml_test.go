package land

import (
	"bytes"
	"testing"
	"unsafe"

	"github.com/ernmw/omwpacker/esm"
	"github.com/stretchr/testify/require"
)

func TestVNMLSize(t *testing.T) {
	require.Equal(t, int(3), int(unsafe.Sizeof(VertexField{})))
}

func TestVNMLFieldUnmarshalMarshal(t *testing.T) {
	const size = vnmlSize

	// Prepare flat input data: sequential bytes 0..(65*65*3-1)
	input := make([]byte, size*size*3)
	for i := range input {
		input[i] = byte(i % 128) // keep in int8 range
	}

	// Create a 65x65 grid of VertexField pointers
	vertices := make([][]VertexField, size)
	for y := range size {
		vertices[y] = make([]VertexField, size)
		for x := range size {
			vertices[y][x] = VertexField{}
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
			offset := (y*size + x) * 3
			if vertex.X != int8(input[offset]) {
				t.Fatalf("vertex[%d][%d].X = %d, want %d", y, x, vertex.X, input[offset])
			}
			if vertex.Y != int8(input[offset+1]) {
				t.Fatalf("vertex[%d][%d].Y = %d, want %d", y, x, vertex.Y, input[offset+1])
			}
			if vertex.Z != int8(input[offset+2]) {
				t.Fatalf("vertex[%d][%d].Z = %d, want %d", y, x, vertex.Z, input[offset+2])
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
