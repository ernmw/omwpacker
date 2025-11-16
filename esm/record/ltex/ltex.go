// LTEX records contain information about landscape textures.
//
//go:generate go run ../generator/gen.go subrecords.json
package ltex

import "github.com/ernmw/omwpacker/esm"

const (
	// LTEX records contain information about landscape textures.
	LTEX esm.RecordTag = "LTEX"
)
