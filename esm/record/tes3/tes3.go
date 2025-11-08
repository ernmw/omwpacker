// Package tes3 handles https://en.uesp.net/wiki/Morrowind_Mod:Mod_File_Format/TES3
package tes3

import (
	"fmt"

	"github.com/ernmw/omwpacker/esm"
)

// NewTES3Record makes a new TES3 record, which must be the first
// record in an ESM.
func NewTES3Record(name string, description string) (*esm.Record, error) {
	// make new empty records
	hedr := &HEDRdata{
		Version:     1.3,
		Flags:       0,
		Name:        name,
		Description: description,
		NumRecords:  0,
	}
	hedrSubRec, err := hedr.Marshal()
	if err != nil {
		return nil, fmt.Errorf("Write HEDR subrecord: %w", err)
	}
	return &esm.Record{
		Tag: TES3,
		Subrecords: []*esm.Subrecord{
			hedrSubRec,
		},
	}, nil
}
