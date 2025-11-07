// Package omwpack bundles up omwscripts files into omwaddon files or vice versa.
// See https://en.uesp.net/wiki/Morrowind_Mod:Mod_File_Format
package omwpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

type pair struct {
	path   string // LUAS
	attach string // LUAF
}

// PackageOmwScripts reads a textual .omwscripts file and writes an .omwaddon file.
// inScriptsPath: path to .omwscripts (text file).
// outAddonPath: path where .omwaddon (binary) will be written.
// templateESPPath: optional path to an existing .esp/.omwaddon to use as a template.
//
//	If provided and that file contains a LUAL record, the function will replace that LUAL
//	record with the newly created one, preserving the template's TES3 header and other data.
//	If empty, the function emits a minimal TES3 header + new LUAL record.
func PackageOmwScripts(inScriptsPath, outAddonPath string) error {
	// 1) Read and parse input scripts text file
	b, err := os.ReadFile(inScriptsPath)
	if err != nil {
		return fmt.Errorf("read scripts file: %w", err)
	}

	lines := strings.Split(string(b), "\n")

	var pairs []pair
	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		// Expect "ATTACH: path"
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line %d: %q (expected 'ATTACH: path')", i+1, line)
		}
		attach := strings.TrimSpace(parts[0])
		path := strings.TrimSpace(parts[1])
		if attach == "" || path == "" {
			return fmt.Errorf("invalid line %d: %q (empty attach or path)", i+1, line)
		}
		pairs = append(pairs, pair{path: path, attach: attach})
	}
	if len(pairs) == 0 {
		return fmt.Errorf("no script pairs found in %s", inScriptsPath)
	}

	// 2) Build the LUAL record bytes
	lual, err := buildLUAL(pairs)
	if err != nil {
		return fmt.Errorf("build LUAL record: %w", err)
	}

	// 3) No template: create minimal TES3 wrapper and put new LUAL inside
	out := bytes.NewBuffer(nil)
	// TES3 header
	out.WriteString("TES3")
	_ = binary.Write(out, binary.LittleEndian, uint32(0)) // placeholder for TES3 size

	// Minimal HEDR subrecord (common in Morrowind TES3):
	// HEDR: version(float32) + numRecords(uint32) + nextObjectID(uint32)
	out.WriteString("HEDR")
	_ = binary.Write(out, binary.LittleEndian, uint32(12))   // size
	_ = binary.Write(out, binary.LittleEndian, float32(1.0)) // version
	_ = binary.Write(out, binary.LittleEndian, uint32(0))    // numRecords (unknown/0)
	_ = binary.Write(out, binary.LittleEndian, uint32(0))    // nextObjectID (0)

	// Append our LUAL
	out.Write(lual)

	// Fill in TES3 size (size of everything after the 8-byte TES3 header)
	outBytes := out.Bytes()
	tes3Size := uint32(len(outBytes) - 8)
	binary.LittleEndian.PutUint32(outBytes[4:], tes3Size)

	if err := os.WriteFile(outAddonPath, outBytes, 0644); err != nil {
		return fmt.Errorf("write out file: %w", err)
	}
	return nil
}

// buildLUAL constructs a LUAL record (bytes) from pairs.
func buildLUAL(pairs []pair) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// LUAL header
	buf.WriteString("LUAL")
	err := binary.Write(buf, binary.LittleEndian, uint32(0)) // placeholder record size
	if err != nil {
		return nil, fmt.Errorf("write LUAL header")
	}

	// reserve spot for size
	_, err = buf.Write(make([]byte, 4))
	if err != nil {
		return nil, fmt.Errorf("write flag padding")
	}

	// Each pair: LUAS (path) then LUAF (attach)
	for _, p := range pairs {
		if err := writeSubrecord(buf, "LUAS", []byte(p.path)); err != nil {
			return nil, fmt.Errorf("write LUAS subrecord")
		}
		if err := writeSubrecord(buf, "LUAF", []byte(p.attach)); err != nil {
			return nil, fmt.Errorf("write LUAF subrecord")
		}
	}

	// patch size
	out := buf.Bytes()
	// size is everything after the 8 byte header (4 id + 4 size), i.e., len(out) - 8
	recSize := uint32(len(out) - 8)
	binary.LittleEndian.PutUint32(out[4:], recSize)
	return out, nil
}

func writeSubrecord(w io.Writer, id string, data []byte) error {
	// id (4 bytes)
	_, err := w.Write([]byte(id))
	if err != nil {
		return fmt.Errorf("write subrecord type")
	}
	// size (uint32)
	err = binary.Write(w, binary.LittleEndian, uint32(len(data)))
	if err != nil {
		return fmt.Errorf("write subrecord length")
	}
	// payload
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("write subrecord data")
	}
	return nil
}

// replaceOrAppendLUAL searches for the first LUAL record in tpl bytes.
// If found, it replaces that LUAL record with newLual and returns (newBytes, true, nil).
// If not found, it appends newLual to the end and returns (newBytes, false, nil).
// This function does minimal parsing: it finds the ASCII "LUAL", reads the uint32 size
// after that, and computes end offset = pos + 8 + size. If parsing fails, we append.
func replaceOrAppendLUAL(tpl, newLual []byte) (out []byte, replaced bool, err error) {
	pos := bytes.Index(tpl, []byte("LUAL"))
	if pos < 0 {
		// not found -> append
		return append([]byte{}, append(tpl, newLual...)...), false, nil
	}
	if pos+8 > len(tpl) {
		// malformed -> append
		return append([]byte{}, append(tpl, newLual...)...), false, nil
	}
	size := binary.LittleEndian.Uint32(tpl[pos+4 : pos+8])
	end := int(pos) + 8 + int(size)
	if end > len(tpl) {
		// malformed -> append
		return append([]byte{}, append(tpl, newLual...)...), false, nil
	}
	// Build new file: everything up to pos, then newLual, then everything after end
	out = make([]byte, 0, len(tpl)-(end-pos)+len(newLual))
	out = append(out, tpl[:pos]...)
	out = append(out, newLual...)
	out = append(out, tpl[end:]...)
	return out, true, nil
}
