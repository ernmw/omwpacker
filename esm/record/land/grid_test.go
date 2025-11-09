package land

import (
	"bytes"
	"testing"
)

func TestFillAndFlattenGridVariableSize(t *testing.T) {
	width, height := 5, 4
	data := make([]byte, width*height)
	for i := range data {
		data[i] = uint8(i)
	}

	// Create a properly sized 2D grid
	grid := make([][]uint8, height)
	for y := range grid {
		grid[y] = make([]uint8, width)
	}

	if err := fillGridFromBytes(grid, width, height, data); err != nil {
		t.Fatalf("FillGridFromBytes failed: %v", err)
	}

	// Verify round-trip correctness
	out := flattenGrid(grid, width, height)
	if !bytes.Equal(out, data) {
		t.Fatalf("flattened data mismatch:\n got %v\nwant %v", out, data)
	}
}

func TestFillGridFromBytesErrors(t *testing.T) {
	width, height := 3, 2
	grid := make([][]uint8, height)
	for i := range grid {
		grid[i] = make([]uint8, width)
	}

	// too short data
	data := make([]byte, width*height-1)
	if err := fillGridFromBytes(grid, width, height, data); err == nil {
		t.Fatalf("expected error for short data, got nil")
	}

	// wrong row width
	grid[0] = make([]uint8, width-1)
	if err := fillGridFromBytes(grid, width, height, make([]byte, width*height)); err == nil {
		t.Fatalf("expected error for wrong row width, got nil")
	}
}
