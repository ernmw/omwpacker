package cell

import (
	"bytes"

	"github.com/ernmw/omwpacker/esm"
)

// These only appear in mod files when creatures or NPCs are moved from one cell to another; they commonly appear in saved game files as things move around.
type MoveReference struct {
	// Reference ID (always the same as the attached FRMR value).
	// Required.
	MVRF *MVRFdata
	// Name of the cell the reference was moved to (interior cells only)
	// zstring
	// Optional.
	CNAM *CNAMdata
	// Coordinates of the cell the reference was moved to (exterior cells only)
	//   int32 - Grid X
	//   int32 - Grid Y
	// Optional.
	CNDT *CNDTdata
	// Reference to the form that was moved.
	// Optional.
	Moved *FormReference
}

type CNAMdata struct {
	zstring
}

func (s *CNAMdata) Tag() esm.SubrecordTag { return "CNAM" }

type zstring string

func (s *zstring) Tag() esm.SubrecordTag { return "???" }

func (s *zstring) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	tmp := zstring(sub.Data)
	s = &tmp
	return nil
}

func (s *zstring) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)
	if _, err := buff.WriteString(string(*s)); err != nil {
		return nil, err
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: buff.Bytes()}, nil
}
