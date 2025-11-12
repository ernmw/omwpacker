package util

import "fmt"

// SliceAsGrid views a 1D slice as a 2D grid of the given width.
func SliceAsGrid[T any](width int, s []T) ([][]T, error) {
	if width <= 0 {
		return nil, fmt.Errorf("width must be positive, got %d", width)
	}

	if len(s)%width != 0 {
		return nil, fmt.Errorf("slice length %d is not a multiple of width %d", len(s), width)
	}

	height := len(s) / width
	grid := make([][]T, height)

	// Manually create slice headers that all point back to the original underlying array (s)
	for i := range height {
		start := i * width
		end := start + width
		grid[i] = s[start:end]
	}

	return grid, nil
}

// GridAsSlice flattens a 2D grid (slice of slices) into a single 1D slice.
// This operation is NOT zero-copy because the inner slices of a [][]T are
// not guaranteed to be contiguous in memory, and typically are not.
func GridAsSlice[T any](grid [][]T) ([]T, error) {
	if len(grid) == 0 {
		return []T{}, nil
	}

	// confirm sizes are the same
	width := len(grid[0])
	for _, row := range grid[1:] {
		if len(row) != width {
			return nil, fmt.Errorf("mismatched inner slices lengths %d and %d", width, len(row))
		}
	}

	// Copy values out.
	result := make([]T, len(grid)*width)
	for i, row := range grid {
		// Determine the start index in the result slice
		start := i * width
		// Copy the row elements into the pre-allocated result slice
		copy(result[start:start+width], row)
	}

	return result, nil
}
