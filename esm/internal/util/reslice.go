package util

import (
	"fmt"
	"unsafe"
)

// SliceFromBytes fills a []T from raw data.
// The returned []T must not change size.
//
// Notes on struct alignment:
// - The offset of each field must be an integer multiple of the field's type size.
// - The total size of the struct must be an integer multiple of the largest field's type size.
func SliceFromBytes[T any](count int, data []uint8) ([]T, error) {
	if count == 0 && len(data) == 0 {
		return make([]T, 0), nil
	}

	var zero T
	elemSize := unsafe.Sizeof(zero)
	if len(data)%int(elemSize) != 0 {
		return nil, fmt.Errorf("byte slice length %d is not a multiple of element size %d", len(data), elemSize)
	}
	elemsNum := len(data) / int(elemSize)
	if count != elemsNum {
		return nil, fmt.Errorf("expected %d elements, got %d", count, elemsNum)
	}

	// Cast the byte slice to a []MyStruct using unsafe.Slice
	// This creates a new slice header pointing to the same underlying memory
	return unsafe.Slice((*T)(unsafe.Pointer(&data[0])), elemsNum), nil
}

// BytesFromSlice fills a []byte from raw data.
// The returned []byte must not change size.
func BytesFromSlice[T any](elems []T) ([]uint8, error) {
	if len(elems) == 0 {
		return make([]uint8, 0), nil
	}
	var zero T
	elemSize := unsafe.Sizeof(zero)
	return unsafe.Slice((*uint8)(unsafe.Pointer(&elems[0])), len(elems)*int(elemSize)), nil
}
