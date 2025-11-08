package cell

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/internal/util"
)

// CellRecord represents a full CellRecord record composed of subrecords.
type CellRecord struct {
	NAME               *NAMEdata
	DATA               *DATAdata
	RGNN               *RGNNdata
	NAM5               *NAM5data
	WHGT               *WHGTdata
	AMBI               *AMBIdata
	MovedReferences    []*MoveReference
	PersistentChildren []*FormReference
	// Count of temporaray children
	NAM0              *NAM0data
	TemporaryChildren []*FormReference
}

// ========== Subrecord: NAME ==========

type NAMEdata struct {
	Name string
}

func (s *NAMEdata) Tag() esm.SubrecordTag { return NAME }

func (s *NAMEdata) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.Name = string(sub.Data)
	return nil
}

func (s *NAMEdata) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)
	if _, err := buff.WriteString(s.Name); err != nil {
		return nil, err
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: buff.Bytes()}, nil
}

// ========== Subrecord: DATA ==========

type DATAdata struct {
	Flags uint32
	GridX int32
	GridY int32
}

func (s *DATAdata) Tag() esm.SubrecordTag { return DATA }

func (s *DATAdata) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	if len(sub.Data) < 12 {
		return fmt.Errorf("CELL.DATA too short: %d < 12", len(sub.Data))
	}
	s.Flags = binary.LittleEndian.Uint32(sub.Data[0:4])
	s.GridX = int32(binary.LittleEndian.Uint32(sub.Data[4:8]))
	s.GridY = int32(binary.LittleEndian.Uint32(sub.Data[8:12]))
	return nil
}

func (s *DATAdata) Marshal() (*esm.Subrecord, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, s.Flags)
	binary.Write(buf, binary.LittleEndian, s.GridX)
	binary.Write(buf, binary.LittleEndian, s.GridY)
	return &esm.Subrecord{Tag: s.Tag(), Data: buf.Bytes()}, nil
}

// ========== Subrecord: RGNN ==========

type RGNNdata struct {
	Value string
}

func (s *RGNNdata) Tag() esm.SubrecordTag { return RGNN }

func (s *RGNNdata) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	s.Value = string(sub.Data)
	return nil
}

func (s *RGNNdata) Marshal() (*esm.Subrecord, error) {
	buff := new(bytes.Buffer)
	if _, err := buff.WriteString(s.Value); err != nil {
		return nil, err
	}
	return &esm.Subrecord{Tag: s.Tag(), Data: buff.Bytes()}, nil
}

// ========== Subrecord: NAM5 ==========

type NAM5data struct {
	Color [3]uint8
}

func (s *NAM5data) Tag() esm.SubrecordTag { return NAM5 }

func (s *NAM5data) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	if len(sub.Data) < 3 {
		return fmt.Errorf("CELL.NAM5 too short: %d < 3", len(sub.Data))
	}
	copy(s.Color[:], sub.Data[:3])
	return nil
}

func (s *NAM5data) Marshal() (*esm.Subrecord, error) {
	return &esm.Subrecord{Tag: s.Tag(), Data: s.Color[:]}, nil
}

// ========== Subrecord: WHGT ==========

type WHGTdata struct {
	WaterHeight float32
}

func (s *WHGTdata) Tag() esm.SubrecordTag { return WHGT }

func (s *WHGTdata) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	if len(sub.Data) < 4 {
		return fmt.Errorf("CELL.WHGT too short")
	}
	s.WaterHeight = util.BytesToFloat32(sub.Data[0:4])
	return nil
}

func (s *WHGTdata) Marshal() (*esm.Subrecord, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, s.WaterHeight)
	return &esm.Subrecord{Tag: s.Tag(), Data: buf.Bytes()}, nil
}

// ========== Subrecord: AMBI ==========

type AMBIdata struct {
	AmbientColor [3]uint8
	Sunlight     [3]uint8
	FogColor     [3]uint8
	FogDensity   float32
}

func (s *AMBIdata) Tag() esm.SubrecordTag { return AMBI }

func (s *AMBIdata) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	if len(sub.Data) < 16 {
		return fmt.Errorf("CELL.AMBI too short: %d < 16", len(sub.Data))
	}
	copy(s.AmbientColor[:], sub.Data[0:3])
	copy(s.Sunlight[:], sub.Data[3:6])
	copy(s.FogColor[:], sub.Data[6:9])
	s.FogDensity = util.BytesToFloat32(sub.Data[12:16])
	return nil
}

func (s *AMBIdata) Marshal() (*esm.Subrecord, error) {
	buf := new(bytes.Buffer)
	buf.Write(s.AmbientColor[:])
	buf.Write(s.Sunlight[:])
	buf.Write(s.FogColor[:])
	buf.Write([]byte{0, 0, 0}) // padding
	binary.Write(buf, binary.LittleEndian, s.FogDensity)
	return &esm.Subrecord{Tag: s.Tag(), Data: buf.Bytes()}, nil
}

// ========== Subrecord: FRMR ==========
// Reference ID (always uint32)
type FRMRdata struct {
	ReferenceID uint32
}

func (s *FRMRdata) Tag() esm.SubrecordTag { return FRMR }

func (s *FRMRdata) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	if len(sub.Data) < 4 {
		return fmt.Errorf("CELL.FRMR too short: %d < 4", len(sub.Data))
	}
	s.ReferenceID = binary.LittleEndian.Uint32(sub.Data[0:4])
	return nil
}

func (s *FRMRdata) Marshal() (*esm.Subrecord, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, s.ReferenceID)
	return &esm.Subrecord{Tag: s.Tag(), Data: buf.Bytes()}, nil
}

// ========== Subrecord: MVRF ==========
// Moved Reference ID (always uint32)
type MVRFdata struct {
	ReferenceID uint32
}

func (s *MVRFdata) Tag() esm.SubrecordTag { return MVRF }

func (s *MVRFdata) Unmarshal(sub *esm.Subrecord) error {
	if s == nil || sub == nil {
		return esm.ErrArgumentNil
	}
	if len(sub.Data) < 4 {
		return fmt.Errorf("CELL.MVRF too short: %d < 4", len(sub.Data))
	}
	s.ReferenceID = binary.LittleEndian.Uint32(sub.Data[0:4])
	return nil
}

func (s *MVRFdata) Marshal() (*esm.Subrecord, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, s.ReferenceID)
	return &esm.Subrecord{Tag: s.Tag(), Data: buf.Bytes()}, nil
}

// ========== CELL Record Parser ==========

// ParseCELL builds a CELL record from a list of subrecords.
func ParseCELL(subs []*esm.Subrecord) (*CellRecord, error) {
	if subs == nil {
		return nil, esm.ErrArgumentNil
	}

	c := &CellRecord{
		MovedReferences:    []*MoveReference{},
		PersistentChildren: []*FormReference{},
		TemporaryChildren:  []*FormReference{},
	}
	for _, sub := range subs {
		switch sub.Tag {
		case NAME:
			s := &NAMEdata{}
			if err := s.Unmarshal(sub); err != nil {
				return nil, err
			}
			c.NAME = s
		case DATA:
			s := &DATAdata{}
			if err := s.Unmarshal(sub); err != nil {
				return nil, err
			}
			c.DATA = s
		case RGNN:
			s := &RGNNdata{}
			if err := s.Unmarshal(sub); err != nil {
				return nil, err
			}
			c.RGNN = s
		case NAM5:
			s := &NAM5data{}
			if err := s.Unmarshal(sub); err != nil {
				return nil, err
			}
			c.NAM5 = s
		case WHGT:
			s := &WHGTdata{}
			if err := s.Unmarshal(sub); err != nil {
				return nil, err
			}
			c.WHGT = s
		case AMBI:
			s := &AMBIdata{}
			if err := s.Unmarshal(sub); err != nil {
				return nil, err
			}
			c.AMBI = s
		case MVRF:
			s := &MVRFdata{}
			if err := s.Unmarshal(sub); err != nil {
				return nil, err
			}
			c.MVRF = append(c.MVRF, s)
		case FRMR:
			s := &FRMRdata{}
			if err := s.Unmarshal(sub); err != nil {
				return nil, err
			}
			c.FRMR = append(c.FRMR, s)
		default:
			return nil, fmt.Errorf("unknown CELL subrecord %q", sub.Tag)
		}
	}
	return c, nil
}
