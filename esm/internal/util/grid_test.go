package util

import (
	"bytes"
	"testing"
	"unsafe"
)

type byteField uint8

func (b *byteField) Data() []byte {
	// Return a 1-byte slice pointing to the underlying ByteField.
	return unsafe.Slice((*byte)(b), 1)
}

func (b *byteField) ByteSize() int {
	return 1
}

func TestZeroAllocGrid_ByteField(t *testing.T) {
	width, height := 3, 2

	// Prepare input data: 0..5
	data := []byte{0, 1, 2, 3, 4, 5}

	// Create grid of ByteField pointers
	grid := make([][]*byteField, height)
	for y := range grid {
		grid[y] = make([]*byteField, width)
		for x := range grid[y] {
			grid[y][x] = new(byteField)
		}
	}

	// Fill grid from data (zero allocation)
	if err := FillGridFromBytes(grid, width, height, data); err != nil {
		t.Fatalf("FillGridFromBytes failed: %v", err)
	}

	// Verify values were copied correctly
	for y := range height {
		for x := range width {
			got := *grid[y][x]      // ByteField
			want := data[y*width+x] // byte
			if byte(got) != want {  // explicit conversion
				t.Fatalf("grid[%d][%d] = %d, want %d", y, x, got, want)
			}
		}
	}

	// Preallocate output buffer for flattening (zero allocation)
	out := make([]byte, width*height)

	// Flatten grid back to out
	if err := FlattenGrid(grid, width, height, out); err != nil {
		t.Fatalf("FlattenGrid failed: %v", err)
	}

	// Check that flattened buffer matches original data
	if !bytes.Equal(out, data) {
		t.Fatalf("flattened data mismatch:\n got  %v\n want %v", out, data)
	}
}
