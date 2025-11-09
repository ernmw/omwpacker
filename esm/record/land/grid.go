package land

import "fmt"

// fillGridFromBytes copies the first (width*height) bytes from data into grid.
// grid must be a slice of rows, each of length width.
// Returns an error if dimensions don't match or data is too short.
func fillGridFromBytes(grid [][]uint8, width, height int, data []byte) error {
	if len(grid) != height {
		return fmt.Errorf("grid height mismatch: got %d, want %d", len(grid), height)
	}
	for i, row := range grid {
		if len(row) != width {
			return fmt.Errorf("grid row %d width mismatch: got %d, want %d", i, len(row), width)
		}
	}
	if len(data) < width*height {
		return fmt.Errorf("not enough data: need %d bytes, got %d", width*height, len(data))
	}

	for y := range height {
		start := y * width
		end := start + width
		copy(grid[y], data[start:end])
	}
	return nil
}

// flattenGrid flattens a 2D grid (height × width) into a contiguous []byte.
func flattenGrid(grid [][]uint8, width, height int) []byte {
	out := make([]byte, width*height)
	for y := range height {
		copy(out[y*width:(y+1)*width], grid[y])
	}
	return out
}
