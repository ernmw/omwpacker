package util_test

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/ernmw/omwpacker/esm/internal/util"
)

// MyStruct must satisfy the constraint ~struct{} and has a known, non-zero size.
// For simplicity and common architecture assumptions, we use fixed-size types.
type MyStruct struct {
	A uint32 // 4 bytes
	B uint32 // 4 bytes
}

const myStructSize = unsafe.Sizeof(MyStruct{})

// Ensure MyStruct is 8 bytes for consistency across tests
func init() {
	if myStructSize != 8 {
		panic("MyStruct size assumption failed, adjust test data")
	}
}

// Helper to create a byte slice representing a MyStruct slice:
// [{A: 1, B: 2}, {A: 3, B: 4}] (assuming little-endian byte ordering)
var testDataBytes = []uint8{
	// Element 1: A=1, B=2
	1, 0, 0, 0, // A=1
	2, 0, 0, 0, // B=2
	// Element 2: A=3, B=4
	3, 0, 0, 0, // A=3
	4, 0, 0, 0, // B=4
}

var testDataStructs = []MyStruct{
	{A: 1, B: 2},
	{A: 3, B: 4},
}

// TestSliceFromBytes validates the conversion from []byte to []T, including error cases.
func TestSliceFromBytes(t *testing.T) {
	// --- Test Cases ---
	tests := []struct {
		name      string
		count     int
		data      []uint8
		expectErr bool
		errMsg    string
		expected  []MyStruct // Only set for successful tests
	}{
		{
			name:      "Success_ValidConversion",
			count:     2,
			data:      testDataBytes,
			expectErr: false,
			expected:  testDataStructs,
		},
		{
			name:      "Error_ByteSliceNotMultipleOfElementSize",
			count:     2,
			data:      testDataBytes[:len(testDataBytes)-1], // 15 bytes total (8*2 - 1)
			expectErr: true,
			errMsg:    "byte slice length 15 is not a multiple of element size 8",
		},
		{
			name:      "Error_CountMismatch",
			count:     1, // Expecting 1 element, but data holds 2
			data:      testDataBytes,
			expectErr: true,
			errMsg:    "expected 1 elements, got 2",
		},
		{
			name:      "Success_EmptySlice",
			count:     0,
			data:      []uint8{},
			expectErr: false,
			expected:  []MyStruct{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := util.SliceFromBytes[MyStruct](tt.count, tt.data)

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
				t.Errorf("SliceFromBytes mismatch:\nExpected: %+v\nGot: %+v", tt.expected, result)
			}
		})
	}
}

// TestBytesFromSlice validates the conversion from []T to []byte.
func TestBytesFromSlice(t *testing.T) {
	// --- Test Cases ---
	tests := []struct {
		name      string
		elems     []MyStruct
		expectErr bool
		expected  []uint8
	}{
		{
			name:      "Success_ValidConversion",
			elems:     testDataStructs,
			expectErr: false,
			expected:  testDataBytes,
		},
		{
			name:      "Success_EmptySlice",
			elems:     []MyStruct{},
			expectErr: false,
			expected:  []uint8{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := util.BytesFromSlice(tt.elems)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("BytesFromSlice mismatch:\nExpected: %v\nGot: %v", tt.expected, result)
			}
		})
	}
}

// TestZeroCopySliceFromBytes verifies that the underlying memory is shared,
// which is the expected behavior of unsafe.Slice for zero-copy casting.
func TestZeroCopySliceFromBytes(t *testing.T) {
	// Use a copy of the bytes so modifications don't leak to other tests
	data := make([]uint8, len(testDataBytes))
	copy(data, testDataBytes)

	// 1. Convert bytes to structs
	structs, err := util.SliceFromBytes[MyStruct](2, data)
	if err != nil {
		t.Fatalf("unexpected error on conversion: %v", err)
	}

	// 2. Modify the struct slice
	structs[0].A = 999 // Change A of the first element

	// Expected byte representation of 999 (0x3E7) in little-endian
	// data[0] should be 0xE7, data[1] should be 0x03, data[2] should be 0x00, data[3] should be 0x00
	expectedBytes := []uint8{
		0xE7, 0x03, 0x00, 0x00, // A=999
		2, 0, 0, 0, // B=2 (unchanged)
		3, 0, 0, 0, // A=3 (unchanged)
		4, 0, 0, 0, // B=4 (unchanged)
	}

	// 3. Verify the change is reflected in the original byte slice
	if !reflect.DeepEqual(data, expectedBytes) {
		t.Errorf("Zero-copy failed: Modifying struct did not update bytes.\nExpected: %v\nGot: %v", expectedBytes, data)
	}
}

// TestZeroCopyBytesFromSlice verifies that the underlying memory is shared.
func TestZeroCopyBytesFromSlice(t *testing.T) {
	// 1. Start with structs
	elems := []MyStruct{
		{A: 10, B: 20},
	}

	// 2. Convert structs to bytes
	data, err := util.BytesFromSlice(elems)
	if err != nil {
		t.Fatalf("unexpected error on conversion: %v", err)
	}

	// 3. Modify the byte slice (change data[4] which corresponds to elems[0].B)
	// B is 20 (0x14). Let's change it to 50 (0x32).
	data[4] = 0x32 // 50 in little-endian

	// 4. Verify the change is reflected in the original struct slice
	expectedStruct := MyStruct{A: 10, B: 50}
	if !reflect.DeepEqual(elems[0], expectedStruct) {
		t.Errorf("Zero-copy failed: Modifying bytes did not update struct.\nExpected: %+v\nGot: %+v", expectedStruct, elems[0])
	}
}
