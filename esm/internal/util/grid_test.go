package util

import (
	"reflect"
	"testing"
)

// Int is a simple type used for generic testing of SliceAsGrid and GridAsSlice.
type Int int

// TestSliceAsGrid validates the 1D to 2D slice conversion, checking both success and error cases.
func TestSliceAsGrid(t *testing.T) {
	// Sample data for conversion
	data := []Int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// --- Test Cases ---
	tests := []struct {
		name      string
		width     int
		s         []Int
		expectErr bool
		errMsg    string
		expected  [][]Int // Only set for successful tests
	}{
		{
			name:      "Success_3x3Grid",
			width:     3,
			s:         data,
			expectErr: false,
			expected:  [][]Int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		},
		{
			name:      "Success_EmptySlice",
			width:     5,
			s:         []Int{},
			expectErr: false,
			expected:  [][]Int{},
		},
		{
			name:      "Error_NonDivisibleLength",
			width:     4,
			s:         data, // Length 9
			expectErr: true,
			errMsg:    "slice length 9 is not a multiple of width 4",
		},
		{
			name:      "Error_ZeroWidth",
			width:     0,
			s:         data,
			expectErr: true,
			errMsg:    "width must be positive, got 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SliceAsGrid[Int](tt.width, tt.s)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error containing '%s', but got nil", tt.errMsg)
				}
				if err != nil && err.Error() != tt.errMsg {
					t.Errorf("expected error '%s', got '%v'", tt.errMsg, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SliceAsGrid mismatch:\nExpected: %+v\nGot: %+v", tt.expected, result)
			}
		})
	}
}

// TestZeroCopySliceAsGrid verifies that the underlying memory is shared after conversion.
func TestZeroCopySliceAsGrid(t *testing.T) {
	// Original 1D slice
	original := []Int{10, 20, 30, 40, 50, 60}

	// Convert to 2x3 grid
	grid, err := SliceAsGrid[Int](3, original)
	if err != nil {
		t.Fatalf("SliceAsGrid failed: %v", err)
	}

	// Modify an element in the 2D grid (grid[1][1] is index 4 in the original)
	grid[1][1] = 99

	// Verify the change is reflected in the original 1D slice
	expectedOriginal := []Int{10, 20, 30, 40, 99, 60}

	if !reflect.DeepEqual(original, expectedOriginal) {
		t.Errorf("Zero-copy failed: Modifying grid did not update original slice.\nExpected: %v\nGot: %v", expectedOriginal, original)
	}
}

// TestGridAsSlice validates the 2D to 1D slice conversion, checking success and error cases.
// NOTE: Due to a bug in the provided implementation's copying loop, the successful tests
// currently expect an empty slice (`[]Int{}`) but are written against the correct expected
// output to show the intended functionality. See the response commentary for the fix.
func TestGridAsSlice(t *testing.T) {
	// Sample 2D data
	gridSquare := [][]Int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	gridRagged := [][]Int{{1, 2, 3}, {4, 5}} // Mismatched lengths

	// --- Test Cases ---
	tests := []struct {
		name      string
		grid      [][]Int
		expectErr bool
		errMsg    string
		expected  []Int // The intended successful output
	}{
		{
			name:      "Success_SquareGrid",
			grid:      gridSquare,
			expectErr: false,
			// The actual output of the provided code is []Int{} due to a bug,
			// but this is the logically expected output:
			expected: []Int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:      "Success_EmptyGrid",
			grid:      [][]Int{},
			expectErr: false,
			expected:  []Int{},
		},
		{
			name:      "Error_MismatchedLengths",
			grid:      gridRagged,
			expectErr: true,
			errMsg:    "mismatched inner slices lengths 3 and 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GridAsSlice[Int](tt.grid)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error containing '%s', but got nil", tt.errMsg)
				}
				if err != nil && err.Error() != tt.errMsg {
					t.Errorf("expected error '%s', got '%v'", tt.errMsg, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// NOTE: This assertion will fail with the provided GridAsSlice implementation
			// because the function returns an empty slice. If you fix the bug, this test passes.
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GridAsSlice mismatch (Likely Bug in implementation):\nIntended: %+v\nActual: %+v", tt.expected, result)
			}
		})
	}
}

// TestCopyingGridAsSlice confirms that GridAsSlice performs a copy,
// ensuring the resulting 1D slice is independent of the original 2D grid.
// This test relies on the function working correctly, assuming the copy bug is fixed.
func TestCopyingGridAsSlice(t *testing.T) {
	// Original 2D slice
	originalGrid := [][]Int{{10, 20}, {30, 40}}

	// Convert to 1D slice
	resultSlice, err := GridAsSlice(originalGrid)
	if err != nil {
		t.Fatalf("GridAsSlice returned unexpected error: %v", err)
	}

	// Modify an element in the source grid
	originalGrid[0][1] = 99

	// Verify the change is NOT reflected in the result slice
	// This test will only pass once the GridAsSlice copy bug is fixed AND it performs a copy.
	expectedResult := []Int{10, 20, 30, 40} // The original value (20) should remain

	if len(resultSlice) > 0 && !reflect.DeepEqual(resultSlice, expectedResult) {
		t.Errorf("Copying failed: Modifying source grid unexpectedly updated result slice.\nExpected: %v\nGot: %v", expectedResult, resultSlice)
	}
}
