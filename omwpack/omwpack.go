// Package omwpack bundles up omwscripts files into omwaddon files or vice versa.
// See https://en.uesp.net/wiki/Morrowind_Mod:Mod_File_Format
package omwpack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type pair struct {
	path   string // LUAS
	attach string // LUAF
}

func writePaddedString(out *bytes.Buffer, s []byte, size int) error {
	if len(s) > size {
		return errors.New("string too big")
	}
	if _, err := out.Write(s); err != nil {
		return err
	}
	if _, err := out.Write(make([]byte, size-len(s))); err != nil {
		return err
	}
	return nil
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
	// https://en.uesp.net/wiki/Morrowind_Mod:Mod_File_Format/TES3
	out.WriteString("HEDR")
	_ = binary.Write(out, binary.LittleEndian, uint32(300))      // size
	_ = binary.Write(out, binary.LittleEndian, float32(1.0))     // version
	_ = binary.Write(out, binary.LittleEndian, uint32(0))        // flags
	_ = writePaddedString(out, []byte("omwpack"), 32)            // company name
	_ = writePaddedString(out, []byte("omwscrips package"), 256) // file desc
	_ = binary.Write(out, binary.LittleEndian, uint32(1))        // num records

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
