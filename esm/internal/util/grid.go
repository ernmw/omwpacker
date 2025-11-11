package util

import (
	"fmt"
)

// BinaryFieldZero represents an element with a slice of underlying bytes.
type BinaryFieldZero interface {
	// Data returns the underlying slice that represents this element
	// in memory. Must be exactly ByteSize() long.
	Data() []byte
	ByteSize() int
}

// FillGridFromBytes fills a 2D grid of BinaryFieldZero elements
// from a flat byte slice. Data is copied directly into each element's backing slice.
func FillGridFromBytes[T BinaryFieldZero](grid [][]T, width, height int, data []byte) error {
	if len(grid) != height {
		return fmt.Errorf("grid height mismatch: got %d, want %d", len(grid), height)
	}
	if height == 0 || width == 0 {
		return nil
	}
	elemSize := grid[0][0].ByteSize()
	if len(data) < width*height*elemSize {
		return fmt.Errorf("not enough data: need %d, got %d", width*height*elemSize, len(data))
	}

	offset := 0
	for y := range height {
		if len(grid[y]) != width {
			return fmt.Errorf("row %d width mismatch: got %d, want %d", y, len(grid[y]), width)
		}
		for x := range width {
			elem := grid[y][x]
			copy(elem.Data(), data[offset:offset+elemSize])
			offset += elemSize
		}
	}
	return nil
}

// FlattenGrid returns a single []byte slice pointing to the grid's elements.
// It requires a preallocated output buffer of size width*height*elemSize.
func FlattenGrid[T BinaryFieldZero](grid [][]T, width, height int, out []byte) error {
	if len(grid) != height {
		return fmt.Errorf("grid height mismatch: got %d, want %d", len(grid), height)
	}
	if height == 0 || width == 0 {
		return nil
	}
	elemSize := grid[0][0].ByteSize()
	if len(out) < width*height*elemSize {
		return fmt.Errorf("output buffer too small: need %d, got %d", width*height*elemSize, len(out))
	}

	offset := 0
	for y := range height {
		if len(grid[y]) != width {
			return fmt.Errorf("row %d width mismatch: got %d, want %d", y, len(grid[y]), width)
		}
		for x := range width {
			elem := grid[y][x]
			copy(out[offset:offset+elemSize], elem.Data())
			offset += elemSize
		}
	}
	return nil
}
